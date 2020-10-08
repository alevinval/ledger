package ledger

import (
	"errors"
	"fmt"
	"os"

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
		writeScanKey      []byte
		isClosed          bool
		messages          chan *Message
		opts              *Options
		db                *storage
		readerChk         *checkpoint
		writerChk         *checkpoint
		writeNotification chan struct{}
		closeNotification chan *Reader
	}

	// Message structure used to represent read results
	Message struct {
		Index uint64
		Data  []byte
	}
)

// NewReader creates a default ledger reader
func (w *Writer) NewReader(id string) (*Reader, error) {
	return w.NewReaderOpts(id, DefaultOptions())
}

// NewReaderOpts creates a customized ledger reader
func (w *Writer) NewReaderOpts(id string, opts *Options) (r *Reader, err error) {
	basePrefix := fmt.Sprintf("ledger-%s-reader-%s", w.id, id)
	logger.Log("reader-prefix", basePrefix)
	r = &Reader{
		id:                id,
		writeScanKey:      buildWriteScanKey(w.basePrefix),
		readerChk:         newCheckpoint(basePrefix, w.db, opts),
		writerChk:         w.chk,
		messages:          make(chan *Message),
		opts:              opts,
		db:                w.db,
		writeNotification: make(chan struct{}, 1),
		closeNotification: w.closedReaderNotification,
	}
	err = r.initialise()
	if err == nil {
		w.newReaderNotification <- r
	}
	return
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
		go r.fetcher()

		// Queue a write notification to ensure there is an initial fetch
		r.writeNotification <- struct{}{}
	}
	return
}

func (r *Reader) Read() <-chan *Message {
	return r.messages
}

func (r *Reader) Close() {
	r.isClosed = true
	r.closeNotification <- r
}

func (r *Reader) fetcher() {
	for !r.isClosed {
		select {
		case <-r.writeNotification:
			r.fetch()
		}
	}
}

func (r *Reader) fetch() {
	cp, err := r.readerChk.GetCheckpoint()
	if err != nil {
		logger.Log("ledger-open", "cannot retrieve reader checkpoint", "error", err)
		return
	}

	startIdx := cp.Index
	logger.Log("ledger-open", "scanning", "prefix", r.writeScanKey, "startIdx", startIdx)
	err = r.db.ScanKeysIndexed(r.writeScanKey, startIdx, func(k []byte, idx uint64) (err error) {
		value, err := r.db.GetBytes(k)
		if err != nil {
			return
		}

		r.messages <- &Message{idx, value}
		r.readerChk.Commit(idx)
		return
	})
}

func buildWriteScanKey(prefix string) []byte {
	return []byte(fmt.Sprintf("%s-write-", prefix))
}
