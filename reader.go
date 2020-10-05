package ledger

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
)

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(os.Stdout)
}

type (
	Reader struct {
		id           string
		isTicking    bool
		writeScanKey []byte
		out          io.Writer
		opts         *Options
		db           *storage
		readerChk    *checkpoint
		writerChk    *checkpoint
		closeCh      chan struct{}
		closedCh     chan struct{}
		mu           *sync.Mutex
	}
)

// NewReader ledger
func NewReader(w *Writer, id string, out io.Writer) (*Reader, error) {
	return NewReaderOpts(w, id, out, DefaultOptions())
}

// NewReaderOpts ledger
func NewReaderOpts(w *Writer, id string, out io.Writer, opts *Options) (*Reader, error) {
	basePrefix := fmt.Sprintf("ledger-%s-reader-%s", w.id, id)
	logger.Log("reader-prefix", basePrefix)
	l := &Reader{
		id:           id,
		writeScanKey: buildWriteScanKey(w.basePrefix),
		readerChk:    newCheckpoint(basePrefix, w.db, opts),
		writerChk:    w.chk,
		out:          out,
		opts:         opts,
		db:           w.db,
		closeCh:      make(chan struct{}),
		closedCh:     make(chan struct{}),
		mu:           new(sync.Mutex),
	}
	return l, l.initialise()
}

func (l *Reader) initialise() (err error) {
	// Ensure a checkpoint exists before anything else
	_, notFoundErr := l.readerChk.GetCheckpoint()
	if notFoundErr != nil {

		// In this case, figure out which offset should be committed
		// based on the configured options (custom, earliest, latest...)
		cp, err := l.readerChk.GetCheckpointFrom(l.writerChk)
		if err == nil {
			err = l.readerChk.Commit(cp.Index)
		} else {
			// This should never happen, if it does the underlying storage
			// may have issues
			logger.Log("error", "cannot create starting checkpoint")
		}
	}
	return
}

func (l *Reader) Open() (err error) {
	cp, err := l.readerChk.GetCheckpoint()
	if err != nil {
		logger.Log("ledger-open", "cannot retrieve reader checkpoint", "error", err)
		return
	}

	startIdx := cp.Index
	logger.Log("ledger-open", "scanning", "prefix", l.writeScanKey, "startIdx", startIdx)
	return l.db.ScanKeysIndexed(l.writeScanKey, startIdx, func(k []byte, idx uint64) (err error) {
		value, err := l.db.GetBytes(k)
		if err != nil {
			return
		}
		_, err = l.out.Write(value)
		if err != nil {
			return
		}
		return l.readerChk.Commit(idx)
	})
}

func (l *Reader) OpenTicker(tickInMs time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.isTicking {
		return
	}

	l.isTicking = true
	l.Open()
	go func() {
		logger.Log("ledger-open-ticker", "started ticking")
		for {
			select {
			case <-timeout(tickInMs):
				l.Open()
			case <-l.closeCh:
				logger.Log("ledger-open-ticker", "stopped ticking")
				l.closedCh <- struct{}{}
				return
			}
		}
	}()
}

func (l *Reader) CloseTicker() {
	l.mu.Lock()
	defer l.mu.Unlock()

	logger.Log("ledger-close", "ledger being closed")
	l.closeCh <- struct{}{}
	<-l.closedCh
	l.isTicking = false
}

func buildWriteScanKey(prefix string) []byte {
	return []byte(fmt.Sprintf("%s-write-", prefix))
}

func timeout(timeoutInMs time.Duration) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		time.Sleep(timeoutInMs * time.Millisecond)
		ch <- struct{}{}
	}()
	return ch
}
