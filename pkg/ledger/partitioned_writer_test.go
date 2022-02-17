package ledger

import (
	"testing"

	"github.com/alevinval/ledger/internal/testutils"
	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/assert"
)

func TestPartitionedWriterWrite(t *testing.T) {
	withPartitionedWriter(t, 4, func(pw *PartitionedWriter) {
		_, err := pw.Write([]byte("some-data"))
		assert.NoError(t, err)
	})
}

func TestPartitionedWriterWrite_returnsWriteOffset(t *testing.T) {
	withPartitionedWriter(t, 3, func(pw *PartitionedWriter) {
		testutils.AssertPartitionedWritesAtOffset(t, pw, 1)
		testutils.AssertPartitionedWritesAtOffset(t, pw, 2)
	})
}

func TestPartitionedWriterWrite_failsWithClosedWriter(t *testing.T) {
	partitions := 3
	withPartitionedWriter(t, partitions, func(pw *PartitionedWriter) {
		pw.Close()

		for i := 0; i < partitions; i++ {
			_, err := pw.Write([]byte("some-data"))
			assert.ErrorIs(t, ErrClosedWriter, err)
		}
	})
}

func withPartitionedWriter(t *testing.T, partitions int, fn func(pw *PartitionedWriter)) {
	testutils.WithDB(func(db *badger.DB) {
		partitionedWriter, err := NewPartitionedWriter("partitioned-writer", db, partitions)
		assert.NoError(t, err)

		defer partitionedWriter.Close()

		fn(partitionedWriter)
	})
}
