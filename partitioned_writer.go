package ledger

import (
	"fmt"

	"github.com/dgraph-io/badger/v2"
)

// PartitionedWriter writes using a partition scheme
type PartitionedWriter struct {
	db         *badger.DB
	writers    []*Writer
	partitions int
	current    int
}

func NewPartitionedWriter(id string, db *badger.DB, partitions int) (*PartitionedWriter, error) {
	writers := make([]*Writer, partitions)
	for i := range writers {
		w, err := NewWriter(fmt.Sprintf("%s-p%d", id, i), db)
		if err != nil {
			return nil, err
		}

		writers[i] = w
	}

	return &PartitionedWriter{
		db:         db,
		writers:    writers,
		partitions: partitions,
		current:    -1,
	}, nil
}

func (w *PartitionedWriter) Write(message []byte) {
	w.writers[w.next()].Write(message)
}

func (w *PartitionedWriter) next() int {
	w.current++
	if w.current >= w.partitions {
		w.current = 0
	}
	return w.current
}
