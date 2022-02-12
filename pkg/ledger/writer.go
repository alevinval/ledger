package ledger

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/alevinval/ledger/internal/checkpoint"
	"github.com/alevinval/ledger/internal/storage"
	"github.com/alevinval/ledger/pkg/proto"
	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
)

var (
	ErrClosedWriter = errors.New("writer is closed")
)

var (
	writeKeyFmt = "%s-write-%0." + strconv.Itoa(storage.UintDigits) + "d"
)

// Writer knows how to store messages in the ledger, it performs
// indexed insertions and keeps track of the index of the latest
// message that has been written. Indexed insertions are possible
// using a BadgerDB Sequence which yields monotonically increasing
// integers.
type Writer struct {
	id             string
	basePrefix     string
	isClosed       bool
	mu             sync.RWMutex
	checkpoint     *checkpoint.Checkpoint
	storage        *storage.Storage
	writeSeq       *badger.Sequence
	writerListener *writerListener
}

// NewWriter creates a default ledger writer
func NewWriter(id string, db *badger.DB) (*Writer, error) {
	return NewWriterOpts(id, db, DefaultOptions())
}

// NewWriterOpts creates a customised ledger Writer
func NewWriterOpts(id string, db *badger.DB, opts *Options) (*Writer, error) {
	basePrefix := fmt.Sprintf("ledger-%s", id)
	writeSeq, err := db.GetSequence(buildWriteSeqKey(basePrefix), opts.SequenceBandwidth)
	if err != nil {
		return nil, err
	}

	storage := &storage.Storage{DB: db, Opts: opts}
	w := &Writer{
		id:         id,
		basePrefix: basePrefix,
		checkpoint: checkpoint.NewCheckpoint(basePrefix, storage, opts),
		storage:    storage,
		writeSeq:   writeSeq,
		writerListener: &writerListener{
			readers:            make(map[string]*Reader),
			newReader:          make(chan *Reader),
			newWrite:           make(chan emptyObj),
			closeManager:       make(chan emptyObj),
			closeManagerNotify: make(chan emptyObj),
		},
	}

	err = w.initialise()
	if err != nil {
		logger.Error("failed writer initialisation", zap.String("id", w.id), zap.Error(err))
		return w, err
	}

	logger.Debug("writer created", zap.String("id", w.id), zap.String("basePrefix", basePrefix))
	return w, err
}

func (w *Writer) initialise() (err error) {
	_, notFoundErr := w.checkpoint.GetCheckpoint()
	if notFoundErr != nil {
		logger.Info("no previous checkpoint found: committing default checkpoint", zap.String("id", w.id))
		err = w.checkpoint.Commit(0)
		if err != nil {
			logger.Error("cannot initialise writer: cannot commit offset", zap.String("id", w.id), zap.Error(err))
			return err
		}
	}
	go w.writerListener.Listen()
	return
}

// Write a payload to the log, updates writer offset.
func (w *Writer) Write(message []byte) (uint64, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if w.isClosed {
		return 0, ErrClosedWriter
	}

	var idx uint64

	idx, err := w.writeSeq.Next()
	if err != nil {
		return 0, err
	}

	// Always skip index zero because our scan is not inclusive and uint64 is used
	// to represent offsets. Because less than zero cannot be represented, zero value
	// is reserved to indicate a scan from earliest available offset.
	if idx == 0 {
		idx, err = w.writeSeq.Next()
		if err != nil {
			return 0, err
		}
	}

	key := buildWriteKey(w.basePrefix, idx)

	err = w.storage.PutBytes(key, message)
	if err != nil {
		return 0, err
	}

	err = w.checkpoint.Commit(idx)
	if err != nil {
		logger.Error("failed committing write", zap.String("id", w.id), zap.Uint64("offset", idx))
		return 0, err
	}

	logger.Debug("committed write", zap.String("id", w.id), zap.Uint64("offset", idx))

	defer w.writerListener.notifyWrite()

	return idx, nil
}

func (w *Writer) GetCheckpoint() (*proto.Checkpoint, error) {
	return w.checkpoint.GetCheckpoint()
}

// Close the writer by releasing the sequence. Not releasing the sequence
// leads to gaps in the number space. A gap that is big enough will break
// the ledger since it relies on fast scans by assuming there are no gaps.
func (w *Writer) Close() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.isClosed {
		return
	}

	w.writerListener.close()
	w.writeSeq.Release()
	w.isClosed = true
}

func buildWriteSeqKey(prefix string) []byte {
	return []byte(fmt.Sprintf("seq-%s-write", prefix))
}

func buildWriteKey(prefix string, seq uint64) []byte {
	return []byte(fmt.Sprintf(writeKeyFmt, prefix, seq))
}
