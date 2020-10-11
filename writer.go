package ledger

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/dgraph-io/badger/v2"
)

// Writer knows how to store messages in the ledger, it performs
// indexed insertions and keeps track of the index of the latest
// message that has been written. Indexed insertions are possible
// using a BadgerDB Sequence which yields monotonically increasing
// integers.
type Writer struct {
	id         string
	basePrefix string
	isClosed   bool
	chk        *checkpoint
	db         *storage
	seq        *badger.Sequence
	listener   *writerListener
	mu         sync.RWMutex
}

// NewWriter creates a default ledger writer
func NewWriter(id string, db *badger.DB) (*Writer, error) {
	return NewWriterOpts(id, db, DefaultOptions())
}

// NewWriterOpts creates a customised ledger Writer
func NewWriterOpts(id string, db *badger.DB, opts *Options) (*Writer, error) {
	basePrefix := fmt.Sprintf("ledger-%s", id)
	logger.Log("writer-prefix", basePrefix)
	seq, err := db.GetSequence(buildWriteSeqKey(basePrefix), opts.SequenceBandwidth)
	if err != nil {
		return nil, err
	}

	s := &storage{db, opts}
	w := &Writer{
		id:         id,
		basePrefix: basePrefix,
		chk:        newCheckpoint(basePrefix, s, opts),
		db:         s,
		seq:        seq,
		listener: &writerListener{
			readers:            make(map[string]*Reader),
			newReader:          make(chan *Reader),
			newWrite:           make(chan struct{}),
			closeManager:       make(chan struct{}),
			closeManagerNotify: make(chan struct{}),
		},
	}

	return w, w.initialise()
}

func (w *Writer) initialise() (err error) {
	_, notFoundErr := w.chk.GetCheckpoint()
	if notFoundErr != nil {
		err = w.chk.Commit(0)
	}
	go w.listener.manager()
	return
}

// Write a payload to the log, updates writer offset.
func (w *Writer) Write(message []byte) (uint64, error) {
	var idx uint64

	idx, err := w.seq.Next()
	if err != nil {
		return 0, err
	}

	// Always skip index zero because our scan is not inclusive and uint64 is used
	// to represent offsets. Because less than zero cannot be represented, zero value
	// is reserved to indicate a scan from earliest available offset.
	if idx == 0 {
		idx, err = w.seq.Next()
		if err != nil {
			return 0, err
		}
	}

	key := buildWriteKey(w.basePrefix, idx)

	err = w.db.PutBytes(key, message)
	if err != nil {
		return 0, err
	}

	defer w.listener.notifyWrite()

	return idx, w.chk.Commit(idx)
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

	w.listener.close()
	w.seq.Release()
	w.isClosed = true

	return
}

func buildWriteSeqKey(prefix string) []byte {
	return []byte(fmt.Sprintf("seq-%s-write", prefix))
}

func buildWriteKey(prefix string, seq uint64) []byte {
	sfmt := "%s-write-%0." + strconv.Itoa(uintDigits) + "d"
	return []byte(fmt.Sprintf(sfmt, prefix, seq))
}
