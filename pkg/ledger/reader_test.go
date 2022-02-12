package ledger

import (
	"testing"
	"time"

	"github.com/alevinval/ledger/internal/testutils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

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
	logs, restore := testutils.CaptureLogs(zap.WarnLevel)
	defer restore()

	withReader(t, func(writer *Writer, reader *Reader) {
		// Set extremely low read shortTimeout
		shortTimeout := 10 * time.Millisecond
		reader.opts.DeliveryTimeout = shortTimeout

		testutils.AssertWrites(t, writer, "one", "two", "three")
		testutils.AssertReads(t, reader, "one", "two", "three")

		testutils.AssertWrites(t, writer, "asd")
		time.Sleep(2 * shortTimeout)
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
