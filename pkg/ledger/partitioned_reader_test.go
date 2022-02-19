package ledger

import (
	"testing"

	"github.com/alevinval/ledger/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestNewPartitionedReader_fromClosedWriter(t *testing.T) {
	withPartitionedWriter(t, 4, func(writer *PartitionedWriter) {
		writer.Close()

		_, err := writer.NewReader("reader")
		assert.Equal(t, err, ErrClosedWriter)
	})
}

func TestPartitionedReader_readWrites(t *testing.T) {
	withPartitionedReader(t, func(writer *PartitionedWriter, reader *PartitionedReader) {
		testutils.AssertWrites(t, writer, "one", "two", "three")

		testutils.AssertReads(t, reader, "one", "two", "three")
	})
}

func TestPartitionedReader_close(t *testing.T) {
	withPartitionedReader(t, func(writer *PartitionedWriter, reader *PartitionedReader) {
		reader.Close()

		_, err := reader.Read()
		assert.ErrorIs(t, err, ErrClosedReader)
	})
}

func withPartitionedReader(t *testing.T, fn func(w *PartitionedWriter, r *PartitionedReader)) {
	withPartitionedWriter(t, 4, func(writer *PartitionedWriter) {
		reader, err := writer.NewReader("reader")
		assert.NoError(t, err)

		defer reader.Close()

		fn(writer, reader)
	})
}
