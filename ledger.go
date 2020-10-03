package ledger

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"

	"github.com/dgraph-io/badger"
	"github.com/go-kit/kit/log"
)

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(os.Stdout)
}

type (
	MasterLedger struct {
		id         string
		basePrefix string
		db         *storage
		seq        *badger.Sequence
		mu         *sync.Mutex
		opts       *Options
		chk        *checkpoint
	}

	Ledger struct {
		id               string
		masterBasePrefix string
		out              io.Writer
		opts             *Options
		db               *storage
		chk              *checkpoint
		chkMaster        *checkpoint
	}

	Options struct {
		Mode              OptMode
		CustomIndex       uint64
		SequenceBandwidth uint64

		// KeySpaceBatchSize must be apower of ten (10, 100, 1000, 10000 ...)
		// When set to something else, the number of digits will determine which
		// power is used.
		KeySpaceBatchSize uint64
	}

	OptMode byte
)

const (
	ModeEarliest OptMode = iota
	ModeLatest
	ModeCustom
)

func DefaultOptions() *Options {
	return &Options{
		Mode:              ModeLatest,
		SequenceBandwidth: 1000,
		KeySpaceBatchSize: 1000,
	}
}

// New slave ledger
func New(master *MasterLedger, id string, out io.Writer) (*Ledger, error) {
	return NewOpts(master, id, out, DefaultOptions())
}

// New slave ledger with custom options
func NewOpts(master *MasterLedger, id string, out io.Writer, opts *Options) (*Ledger, error) {
	basePrefix := fmt.Sprintf("ledger-%s-slave-%s", master.id, id)
	logger.Log("slave-prefix", basePrefix)
	l := &Ledger{
		id:               id,
		masterBasePrefix: master.basePrefix,
		chk:              newCheckpoint(basePrefix, master.db, opts),
		chkMaster:        master.chk,
		out:              out,
		opts:             opts,
		db:               master.db,
	}
	return l, l.initialise()
}

func NewMaster(id string, db *badger.DB) (*MasterLedger, error) {
	return NewMasterOpts(id, db, DefaultOptions())
}

// NewMaster ledger creation
func NewMasterOpts(id string, db *badger.DB, opts *Options) (*MasterLedger, error) {
	basePrefix := fmt.Sprintf("ledger-%s", id)
	logger.Log("master-prefix", basePrefix)
	seq, err := db.GetSequence(buildWriteSeqKey(basePrefix), opts.SequenceBandwidth)
	if err != nil {
		return nil, err
	}

	l := &MasterLedger{
		id:         id,
		basePrefix: basePrefix,
		db:         &storage{db, opts},
		seq:        seq,
		mu:         new(sync.Mutex),
		opts:       opts,
		chk:        newCheckpoint(basePrefix, &storage{db, opts}, opts),
	}

	return l, l.initialise()
}

func (l *MasterLedger) initialise() (err error) {
	// Ensure a checkpoint exists before anything else
	_, notFoundErr := l.chk.GetCheckpoint()
	if notFoundErr != nil {
		err = l.chk.Commit(0)
	}
	return
}

func (l *Ledger) initialise() (err error) {
	// Ensure a checkpoint exists before anything else
	_, notFoundErr := l.chk.GetCheckpoint()
	if notFoundErr != nil {

		// In this case, figure out which offset should be committed
		// based on the configured options (custom, earliest, latest...)
		cp, err := l.chk.GetCheckpointFrom(l.chkMaster)
		if err == nil {
			err = l.chk.Commit(cp.Index)
		} else {
			// This should never happen, if it does the underlying storage
			// may have issues
			logger.Log("error", "cannot create starting checkpoint")
		}
	}
	return
}

func (l *MasterLedger) Write(message []byte) error {
	var idx uint64

	idx, err := l.seq.Next()
	if err != nil {
		return err
	}

	// Always skip index zero because our scan is not inclusive and uint64 is used
	// to represent offsets. Because less than zero cannot be represented, zero value
	// is reserved to indicate a scan from earliest available offset.
	if idx == 0 {
		idx, err = l.seq.Next()
		if err != nil {
			return err
		}
	}

	key := buildWriteKey(l.basePrefix, idx)
	logger.Log("writeKey", key)

	err = l.db.PutBytes(key, message)
	if err != nil {
		return err
	}

	return l.chk.Commit(idx)
}

// Close the ledger by releasing the sequence. Not releasing the sequence
// leads to gaps in the number space. A gap big enough will break the ledger
// since it relies on fast scans by assuming there are no gaps (once a gap
// is detected, scanning is aborted, since EOF is assumed)
func (l *MasterLedger) Close() {
	l.seq.Release()
}

func (l *Ledger) Open() (err error) {
	cp, err := l.chk.GetCheckpoint()
	if err != nil {
		logger.Log("ledger-open", "cannot retrieve slave checkpoint", "error", err)
		return
	}

	startIdx := cp.Index
	scanPrefix := l.buildWriteScanKey()
	logger.Log("ledger-open", "scanning", "prefix", scanPrefix, "startIdx", startIdx)
	return l.db.ScanKeysIndexed(scanPrefix, startIdx, func(k []byte, idx uint64) (err error) {
		logger.Log("ledger-open", "replaying", "key", k)
		value, err := l.db.GetBytes(k)
		if err != nil {
			return
		}
		_, err = l.out.Write(value)
		if err != nil {
			return
		}
		return l.chk.Commit(idx)
	})
}

func (l *Ledger) buildWriteScanKey() []byte {
	return []byte(fmt.Sprintf("%s-write", l.masterBasePrefix))
}

func buildWriteSeqKey(prefix string) []byte {
	return []byte(fmt.Sprintf("seq-%s-write", prefix))
}

func buildWriteKey(prefix string, seq uint64) []byte {
	sfmt := "%s-write-%0." + strconv.Itoa(uintDigits) + "d"
	return []byte(fmt.Sprintf(sfmt, prefix, seq))
}
