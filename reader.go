package ledger

import (
	"errors"
	"fmt"
	"os"
	"time"

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
		writeScanKey []byte
		isClosed     bool
		messages     chan *Message
		opts         *Options
		db           *storage
		readerChk    *checkpoint
		writerChk    *checkpoint
	}

	// Message structure used to represent read results
	Message struct {
		Index uint64
		Data  []byte
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
		writeScanKey: buildWriteScanKey(w.basePrefix),
		readerChk:    newCheckpoint(basePrefix, w.db, opts),
		writerChk:    w.chk,
		messages:     make(chan *Message),
		opts:         opts,
		db:           w.db,
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
	if err == nil {
		go r.fetcher(r.opts.FetchIntervalMs)
	}
	return
}

func (r *Reader) Read() <-chan *Message {
	return r.messages
}

func (r *Reader) Close() {
	r.isClosed = true
}

func (r *Reader) fetcher(emptyFetchIntervalMs int) {
	for !r.isClosed {
		results, err := r.fetch()
		if err != nil {
			logger.Log("ledger-fetch", "error fetching", "error", err)
			time.Sleep(time.Duration(emptyFetchIntervalMs) * time.Millisecond)
			continue
		}
		if len(results) == 0 {
			logger.Log("ledger-fetch", "empty fetch")
			time.Sleep(time.Duration(emptyFetchIntervalMs) * time.Millisecond)
			continue
		}
		for i := range results {
			r.messages <- results[i]
			r.readerChk.Commit(results[i].Index)
		}
	}
}

func (r *Reader) fetch() ([]*Message, error) {
	cp, err := r.readerChk.GetCheckpoint()
	if err != nil {
		logger.Log("ledger-open", "cannot retrieve reader checkpoint", "error", err)
		return nil, err
	}

	results := make([]*Message, 0, r.opts.BatchSize)
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

		results = append(results, &Message{idx, value})
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
