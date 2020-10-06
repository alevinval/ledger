package ledger

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/go-kit/kit/log"
)

var (
	logger      log.Logger
	errFullPool = errors.New("results pool is full")
)

func init() {
	logger = log.NewLogfmtLogger(os.Stdout)
}

type (
	Reader struct {
		id                string
		isTicking         bool
		writeScanKey      []byte
		currentMessageIdx int
		currentMessage    *message
		messagePool       []message
		opts              *Options
		db                *storage
		readerChk         *checkpoint
		writerChk         *checkpoint
		closeCh           chan struct{}
		closedCh          chan struct{}
		mu                *sync.Mutex
	}

	// intermediate structure used to buffer reads
	message struct {
		idx  uint64
		data []byte
	}
)

// NewReader ledger
func NewReader(w *Writer, id string) (*Reader, error) {
	return NewReaderOpts(w, id, DefaultOptions())
}

// NewReaderOpts ledger
func NewReaderOpts(w *Writer, id string, opts *Options) (*Reader, error) {
	basePrefix := fmt.Sprintf("ledger-%s-reader-%s", w.id, id)
	logger.Log("reader-prefix", basePrefix)
	l := &Reader{
		id:                id,
		writeScanKey:      buildWriteScanKey(w.basePrefix),
		readerChk:         newCheckpoint(basePrefix, w.db, opts),
		writerChk:         w.chk,
		currentMessageIdx: 0,
		messagePool:       make([]message, 0),
		opts:              opts,
		db:                w.db,
		closeCh:           make(chan struct{}),
		closedCh:          make(chan struct{}),
		mu:                new(sync.Mutex),
	}
	return l, l.initialise()
}

func (r *Reader) initialise() (err error) {
	// Ensure a checkpoint exists before anything else
	_, notFoundErr := r.readerChk.GetCheckpoint()
	if notFoundErr != nil {

		// In this case, figure out which offset should be committed
		// based on the configured options (custom, earliest, latest...)
		cp, err := r.readerChk.GetCheckpointFrom(r.writerChk)
		if err == nil {
			err = r.readerChk.Commit(cp.Index)
		} else {
			// This should never happen, if it does the underlying storage
			// may have issues
			logger.Log("ledger-initialise", "cannot create starting checkpoint")
		}
	}
	return
}

func (r *Reader) Read(out []byte) (n int, err error) {
	if r.currentMessage == nil {
		if r.currentMessageIdx >= len(r.messagePool) {
			r.currentMessageIdx = 0
			r.messagePool, err = r.fetch()
			if err != nil {
				return 0, err
			}
			if len(r.messagePool) == 0 {
				return 0, io.EOF
			}
		}
		r.currentMessage = &r.messagePool[r.currentMessageIdx]
		r.currentMessageIdx++
	}

	if len(out) < len(r.currentMessage.data) {
		return 0, io.ErrShortBuffer
	}

	var l = len(r.currentMessage.data)
	copy(out, r.currentMessage.data)
	r.readerChk.Commit(r.currentMessage.idx)

	r.currentMessage = nil
	return l, io.EOF
}

func (r *Reader) fetch() ([]message, error) {
	cp, err := r.readerChk.GetCheckpoint()
	if err != nil {
		logger.Log("ledger-open", "cannot retrieve reader checkpoint", "error", err)
		return nil, err
	}

	results := make([]message, 0, r.opts.BatchSize)
	startIdx := cp.Index
	logger.Log("ledger-open", "scanning", "prefix", r.writeScanKey, "startIdx", startIdx)
	err = r.db.ScanKeysIndexed(r.writeScanKey, startIdx, func(k []byte, idx uint64) (err error) {
		if len(results) >= cap(results) {
			return errFullPool
		}

		value, err := r.db.GetBytes(k)
		if err != nil {
			return
		}

		results = append(results, message{idx, value})
		return
	})
	if err == errFullPool {
		return results, nil
	}
	return results, err
}

func buildWriteScanKey(prefix string) []byte {
	return []byte(fmt.Sprintf("%s-write-", prefix))
}
