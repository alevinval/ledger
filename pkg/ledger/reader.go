package ledger

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/alevinval/ledger/internal/checkpoint"
	"github.com/alevinval/ledger/internal/util/log"
	"go.uber.org/zap"
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

		messages           chan *Message
		triggerFetch       chan emptyObj
		fetcherClose       chan emptyObj
		fetcherCloseNotify chan emptyObj
	}

	// Message structure used to represent read results
	Message struct {
		Offset uint64
		Data   []byte
	}

	// PartitionedMessage structure that represents messages from a partitioned reader
	PartitionedMessage struct {
		Partition Commitable
		Offset    uint64
		Data      []byte
	}

	// Commitable interface
	Commitable interface {
		Commit(uint64) error
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
		return nil, io.ErrClosedPipe
	}

	basePrefix := fmt.Sprintf("ledger-%s-reader-%s", w.id, id)
	r := &Reader{
		id:           id,
		writeScanKey: buildWriteScanKey(w.basePrefix),
		checkpoint:   checkpoint.NewCheckpoint(basePrefix, w.storage, opts),
		opts:         opts,
		writer:       w,

		messages:           make(chan *Message),
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
	// Ensure a checkpoint exists before anything else
	_, notFoundErr := r.checkpoint.GetCheckpoint()
	if notFoundErr != nil {

		// In this case, figure out which offset should be committed
		// based on the configured options (custom, earliest, latest...)
		cp, err := r.checkpoint.GetCheckpointFrom(r.writer.checkpoint)
		if err == nil {
			err = r.checkpoint.Commit(cp.Offset)
			if err != nil {
				logger.Error("cannot initialise reader: cannot commit offset", zap.Error(err))
				return err
			}
		} else {
			// This should never happen, if it does the underlying storage
			// may have issues
			logger.Error("cannot initialise reader: cannot create starting checkpoint", zap.Error(err))
			return err
		}
	}
	go r.fetcher()
	r.doTriggerFetch()

	return nil
}

// Read the ledger, returns a channel where messages can be received
// by the consumer of this API.
// An error is returned in case the Reader is already closed.
func (r *Reader) Read() (<-chan *Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.isClosed {
		return nil, io.ErrClosedPipe
	}

	return r.messages, nil
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
			return
		}

		select {
		case r.messages <- &Message{offset, value}:
			r.fetchStartOffset = offset
		case <-time.After(r.opts.DeliveryTimeout * time.Millisecond):
			logger.Warn("reader delivery timeout: make sure messages are being consumed", zap.String("id", r.id))
			return
		}
		return
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
