package ledger

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/go-kit/kit/log"
)

var (
	logger log.Logger
)

func init() {
	logger = log.NewLogfmtLogger(os.Stdout)
}

type (
	Reader struct {
		id                 string
		writeScanKey       []byte
		isClosed           bool
		messages           chan *Message
		opts               *Options
		chk                *checkpoint
		w                  *Writer
		triggerFetch       chan struct{}
		fetcherClose       chan struct{}
		fetcherCloseNotify chan struct{}
		mu                 sync.RWMutex
	}

	// Message structure used to represent read results
	Message struct {
		Offset uint64
		Data   []byte
	}
)

// NewReader creates a default ledger reader
func (w *Writer) NewReader(id string) (*Reader, error) {
	return w.NewReaderOpts(id, DefaultOptions())
}

// NewReaderOpts creates a customized ledger reader
func (w *Writer) NewReaderOpts(id string, opts *Options) (r *Reader, err error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if w.isClosed {
		return nil, io.ErrClosedPipe
	}

	basePrefix := fmt.Sprintf("ledger-%s-reader-%s", w.id, id)
	logger.Log("reader-prefix", basePrefix)
	r = &Reader{
		id:                 id,
		writeScanKey:       buildWriteScanKey(w.basePrefix),
		chk:                newCheckpoint(basePrefix, w.db, opts),
		messages:           make(chan *Message),
		opts:               opts,
		w:                  w,
		triggerFetch:       make(chan struct{}),
		fetcherClose:       make(chan struct{}),
		fetcherCloseNotify: make(chan struct{}),
	}
	err = r.initialise()
	if err == nil {
		w.listener.notifyReader(r)
	}
	return
}

func (r *Reader) initialise() (err error) {
	// Ensure a checkpoint exists before anything else
	_, notFoundErr := r.chk.GetCheckpoint()
	if notFoundErr != nil {

		// In this case, figure out which offset should be committed
		// based on the configured options (custom, earliest, latest...)
		cp, err := r.chk.GetCheckpointFrom(r.w.chk)
		if err == nil {
			err = r.chk.Commit(cp.Offset)
		} else {
			// This should never happen, if it does the underlying storage
			// may have issues
			logger.Log("ledger-initialise", "cannot create starting checkpoint")
		}
	}
	if err == nil {
		go r.fetcher()
		r.doTriggerFetch()
	}
	return
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

// Close the reader and stop fetching records.
func (r *Reader) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.isClosed {
		return
	}

	r.fetcherClose <- struct{}{}
	<-r.fetcherCloseNotify
	close(r.messages)
	r.isClosed = true
}

func (r *Reader) doTriggerFetch() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.isClosed {
		return true
	}

	r.triggerFetch <- struct{}{}
	return false
}

func (r *Reader) fetch() {
	cp, err := r.chk.GetCheckpoint()
	if err != nil {
		logger.Log("ledger-open", "cannot retrieve reader checkpoint", "error", err)
		return
	}

	startOffset := cp.Offset
	logger.Log("ledger-open", "scanning", "prefix", r.writeScanKey, "startOffset", startOffset)
	err = r.w.db.ScanKeysIndexed(r.writeScanKey, startOffset, func(k []byte, offset uint64) (err error) {
		value, err := r.w.db.GetBytes(k)
		if err != nil {
			return
		}

		r.messages <- &Message{offset, value}
		r.chk.Commit(offset)
		return
	})
}

func (r *Reader) fetcher() {
	for {
		select {
		case <-r.triggerFetch:
			r.fetch()
		case <-r.fetcherClose:
			r.fetcherCloseNotify <- struct{}{}
			return
		}
	}
}

func buildWriteScanKey(prefix string) []byte {
	return []byte(fmt.Sprintf("%s-write-", prefix))
}
