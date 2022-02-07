package ledger

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
)

// PartitionedWriter writes using a partition scheme
type PartitionedWriter struct {
	id         string
	partitions int
	current    int
	writers    []*Writer
	db         *badger.DB
}

func NewPartitionedWriter(id string, db *badger.DB, partitions int) (*PartitionedWriter, error) {
	writers := make([]*Writer, partitions)
	for i := range writers {
		w, err := NewWriter(fmt.Sprintf("%s-part-%d", id, i), db)
		if err != nil {
			logger.Error("cannot create partitioned writer", zap.Error(err))
			return nil, err
		}

		writers[i] = w
	}

	pw := &PartitionedWriter{
		id:         id,
		partitions: partitions,
		current:    -1,
		db:         db,
		writers:    writers,
	}

	pw.logInfo()

	return pw, nil
}

func (w *PartitionedWriter) Write(message []byte) (uint64, error) {
	return w.writers[w.next()].Write(message)
}

func (w *PartitionedWriter) next() int {
	w.current++
	if w.current >= w.partitions {
		w.current = 0
	}
	return w.current
}

func (pw *PartitionedWriter) logInfo() {
	for _, w := range pw.writers {
		checkpoint, err := w.checkpoint.GetCheckpoint()
		if err != nil {
			logger.Debug("error getting checkpoint", zap.String("partitioned-writer", pw.id), zap.String("writer", w.id), zap.Error(err))
		} else {
			logger.Debug("partitioned writer", zap.String("id", pw.id), zap.String("writer", w.id), zap.Uint64("offset", checkpoint.GetOffset()))
		}
	}
}

func (pw *PartitionedWriter) createReaders(id string) ([]*Reader, error) {
	readers := make([]*Reader, pw.partitions)
	for i := range pw.writers {
		r, err := pw.writers[i].NewReader(fmt.Sprintf("%s-part-%d", id, i))
		if err != nil {
			return nil, err
		}
		readers[i] = r
	}
	return readers, nil
}
