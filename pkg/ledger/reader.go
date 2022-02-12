package ledger

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/alevinval/ledger/internal/base"
	"github.com/alevinval/ledger/internal/checkpoint"
	"github.com/alevinval/ledger/internal/log"
	"github.com/alevinval/ledger/pkg/proto"
	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
)

var (
	// ErrClosedReader is returned when attempting to read from a closed Reader
	ErrClosedReader = errors.New("reader is closed")
)

var (
	logger = log.GetLogger()
)

type (
	Reader struct {
		id               string
		isClosed         bool
		fetchStartOffset uint64
		writeScanKey     []byte
		mu               sync.RWMutex
		opts             *Options
		checkpoint       *checkpoint.Checkpoint
		writer           *Writer

		messages           chan base.Message
		triggerFetch       chan emptyObj
		fetcherClose       chan emptyObj
		fetcherCloseNotify chan emptyObj
	}
)

// NewReader creates a default ledger reader
func (w *Writer) NewReader(id string) (*Reader, error) {
	return w.NewReaderOpts(id, DefaultOptions())
}

// NewReaderOpts creates a customized ledger reader
func (w *Writer) NewReaderOpts(id string, opts *Options) (*Reader, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if w.isClosed {
		return nil, ErrClosedWriter
	}

	basePrefix := fmt.Sprintf("ledger-%s-reader-%s", w.id, id)
	r := &Reader{
		id:           id,
		writeScanKey: buildWriteScanKey(w.basePrefix),
		checkpoint:   checkpoint.NewCheckpoint(basePrefix, w.storage, opts),
		opts:         opts,
		writer:       w,

		messages:           make(chan base.Message),
		triggerFetch:       make(chan emptyObj, 1),
		fetcherClose:       make(chan emptyObj),
		fetcherCloseNotify: make(chan emptyObj),
	}

	err := r.initialise()
	if err != nil {
		logger.Error("failed reader initialisation", zap.String("id", r.id), zap.Error(err))
		return r, err
	}

	w.writerListener.notifyReader(r)

	logger.Info("reader created", zap.String("id", id), zap.String("basePrefix", basePrefix))
	return r, nil
}

func (r *Reader) initialise() error {
	_, err := r.checkpoint.GetCheckpoint()
	if err == badger.ErrKeyNotFound {
		err = r.commitInitialCheckpoint()
		if err != nil {
			logger.Error("cannot initialise reader", zap.Error(err))
			return err
		}
	} else {
		logger.Error("cannot initialise reader: cannot get current checkpoint", zap.Error(err))
	}

	go r.fetcher()

	return nil
}

// Read the ledger, returns a channel where messages can be received
// by the consumer of this API.
// An error is returned in case the Reader is already closed.
func (r *Reader) Read() (<-chan base.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.isClosed {
		return nil, ErrClosedReader
	}

	r.doTriggerFetch()
	return r.messages, nil
}

// GetCheckpoint returns the current checkpoint of the reader
func (r *Reader) GetCheckpoint() (*proto.Checkpoint, error) {
	return r.checkpoint.GetCheckpoint()
}

// Commit the offset, creates a new checkpoint to persist
// lastest ID that has been processed by the Reader.
func (r *Reader) Commit(offset uint64) error {
	return r.checkpoint.Commit(offset)
}

// Close the reader and stop fetching records.
func (r *Reader) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.isClosed {
		return
	}

	fireAndWait(r.fetcherClose, r.fetcherCloseNotify)
	close(r.messages)
	r.isClosed = true
}

// commitInitialOffset commits the reader offset for the first time
func (r *Reader) commitInitialCheckpoint() error {
	initialCheckpoint, err := r.checkpoint.GetCheckpointFrom(r.writer.checkpoint)

	if err != nil {
		logger.Error("cannot commit initial checkpoint", zap.Error(err))
		return err
	}

	err = r.checkpoint.Commit(initialCheckpoint.Offset)
	if err != nil {
		logger.Error("cannot commit initial checkpoint", zap.Error(err))
		return err
	}

	return nil
}

func (r *Reader) doTriggerFetch() (isClosed bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.isClosed {
		return true
	}

	fireAndForget(r.triggerFetch)
	return false
}

func (r *Reader) fetch() {
	if r.fetchStartOffset == 0 {
		cp, err := r.checkpoint.GetCheckpoint()
		if err != nil {
			logger.Error("reader fetch failed: cannot retrieve checkpoint", zap.String("id", r.id), zap.Error(err))
			return
		}
		r.fetchStartOffset = cp.Offset
	}

	logger.Debug("reader fetch scan", zap.String("id", r.id), zap.ByteString("writeScanKey", r.writeScanKey), zap.Uint64("offset", r.fetchStartOffset))
	r.writer.storage.ScanKeysIndexed(r.writeScanKey, r.fetchStartOffset, func(k []byte, offset uint64) (err error) {
		value, err := r.writer.storage.GetBytes(k)
		if err != nil {
			logger.Error("failed scanning keys", zap.String("id", r.id), zap.Error(err))
			return
		}

		select {
		case r.messages <- &messageImpl{offset, value}:
			r.fetchStartOffset = offset
			return
		case <-time.After(r.opts.DeliveryTimeout):
			logger.Warn("reader delivery timeout: make sure messages are being consumed", zap.String("id", r.id))
			return
		}
	})
}

func (r *Reader) fetcher() {
	for {
		select {
		case <-r.triggerFetch:
			r.fetch()
		case <-r.fetcherClose:
			r.fetcherCloseNotify <- empty
			return
		}
	}
}

func buildWriteScanKey(prefix string) []byte {
	return []byte(fmt.Sprintf("%s-write-", prefix))
}
