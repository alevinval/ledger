package ledger

import (
	"testing"
	"time"

	"github.com/alevinval/ledger/internal/log"
	"github.com/alevinval/ledger/internal/testutils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func captureLogs(level zapcore.Level) (logs *observer.ObservedLogs, restore func()) {
	originalLogger := logger.GetZapLogger()
	restoreOriginalLogger := func() {
		log.SetZapLogger(originalLogger)
	}

	core, logs := observer.New(zap.WarnLevel)
	newLogger := zap.New(core)
	log.SetZapLogger(newLogger)
	return logs, restoreOriginalLogger
}

func TestNewReader_fromClosedWriter(t *testing.T) {
	withWriter(t, func(writer *Writer) {
		writer.Close()

		_, err := writer.NewReader("reader")
		assert.Equal(t, err, ErrClosedWriter)
	})
}

func TestReader_readNoWrites(t *testing.T) {
	withReader(t, func(writer *Writer, reader *Reader) {
		testutils.AssertCheckpointAt(t, reader, 0)
		testutils.AssertReads(t, reader, "")
		testutils.AssertCheckpointAt(t, reader, 0)
	})
}

func TestReader_readWrites(t *testing.T) {
	withReader(t, func(writer *Writer, reader *Reader) {
		testutils.AssertWrites(t, writer, "one", "two", "three")

		testutils.AssertCheckpointAt(t, reader, 0)
		testutils.AssertReads(t, reader, "one", "two", "three")
		testutils.AssertCheckpointAt(t, reader, 3)
		println("hey")
	})
}

func TestReader_readWritesWithoutAutoCommit(t *testing.T) {
	withReader(t, func(writer *Writer, reader *Reader) {
		testutils.AssertWrites(t, writer, "one", "two", "three")

		testutils.AssertCheckpointAt(t, reader, 0)
		testutils.AssertReadsNoCommit(t, reader, "one", "two", "three")
		testutils.AssertCheckpointAt(t, reader, 0)
	})
}

func TestReader_slowReads(t *testing.T) {
	logs, restore := captureLogs(zap.WarnLevel)
	defer restore()

	withReader(t, func(writer *Writer, reader *Reader) {
		reader.opts.DeliveryTimeout = 10

		testutils.AssertWrites(t, writer, "one", "two", "three")

		time.AfterFunc(2*reader.opts.DeliveryTimeout*time.Millisecond, func() {
			testutils.AssertWrites(t, writer, "asd")
		})

		testutils.AssertReads(t, reader, "one", "two", "three")

		time.Sleep(250 * time.Millisecond)
		testutils.AssertReads(t, reader, "asd")

		assert.Greater(t, logs.Len(), 0)
		deliveryTimeoutWarn := logs.All()[0]
		assert.Equal(t, "reader delivery timeout: make sure messages are being consumed", deliveryTimeoutWarn.Message)
		assert.Equal(t, "reader", deliveryTimeoutWarn.ContextMap()["id"])
	})
}

func TestReader_close(t *testing.T) {
	withReader(t, func(writer *Writer, reader *Reader) {
		reader.Close()

		_, err := reader.Read()
		assert.ErrorIs(t, err, ErrClosedReader)
	})
}

func withReader(t *testing.T, fn func(w *Writer, r *Reader)) {
	withWriter(t, func(writer *Writer) {
		reader, err := writer.NewReader("reader")
		defer reader.Close()

		assert.NoError(t, err)

		fn(writer, reader)
	})
}
