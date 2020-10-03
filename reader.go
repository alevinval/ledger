package ledger

import (
	"fmt"
	"io"
	"os"

	"github.com/go-kit/kit/log"
)

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(os.Stdout)
}

type (
	Reader struct {
		id           string
		writeScanKey []byte
		out          io.Writer
		opts         *Options
		db           *storage
		readerChk    *checkpoint
		writerChk    *checkpoint
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
		logger.Log("ledger-open", "replaying", "key", k)
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

func (l *Reader) Commit(idx uint64) error {
	return l.readerChk.Commit(idx)
}

func buildWriteScanKey(prefix string) []byte {
	return []byte(fmt.Sprintf("%s-write", prefix))
}
