package ledger

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"

	"github.com/alevinval/ledger/pkg/proto"
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
	}

	Ledger struct {
		basePrefix string
		out        io.Writer
		master     *MasterLedger
		opts       *Options
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
func New(master *MasterLedger, id string, out io.Writer) *Ledger {
	return NewOpts(master, id, out, DefaultOptions())
}

// New slave ledger with custom options
func NewOpts(master *MasterLedger, id string, out io.Writer, opts *Options) *Ledger {
	basePrefix := fmt.Sprintf("ledger-%s-slave-%s", master.id, id)
	logger.Log("slave-prefix", basePrefix)
	l := &Ledger{
		basePrefix: basePrefix,
		out:        out,
		master:     master,
		opts:       opts,
	}
	_ = l.master.getCheckPoint(l.basePrefix, opts)
	return l
}

func NewMaster(id string, db *badger.DB) *MasterLedger {
	return NewMasterOpts(id, db, DefaultOptions())
}

// NewMaster ledger creation
func NewMasterOpts(id string, db *badger.DB, opts *Options) *MasterLedger {
	basePrefix := fmt.Sprintf("ledger-%s", id)
	logger.Log("master-prefix", basePrefix)
	seq, err := db.GetSequence(buildWriteSeqKey(basePrefix), opts.SequenceBandwidth)
	if err != nil {
		logger.Log("get_sequence", err)
		panic("could not retrieve a sequence from badger")
	}

	return &MasterLedger{
		id:         id,
		basePrefix: basePrefix,
		db:         &storage{db, opts},
		seq:        seq,
		mu:         new(sync.Mutex),
		opts:       opts,
	}
}

func (l *MasterLedger) Write(message []byte) {
	var idx uint64
	var err error

	for {
		idx, err = l.seq.Next()
		if err == nil && idx != 0 {
			break
		}
	}

	key := buildWriteKey(l.basePrefix, idx)
	logger.Log("writeKey", key)
	l.db.PutBytes(key, message)
	l.commit(l.basePrefix, idx)
}

// Close the ledger by releasing the sequence. Not releasing the sequence
// leads to gaps in the number space. A gap big enough will break the ledger
// since it relies on fast scans by assuming there are no gaps (once a gap
// is detected, scanning is aborted, since EOF is assumed)
func (l *MasterLedger) Close() {
	l.seq.Release()
}

func (l *MasterLedger) getCheckPoint(prefix string, opts *Options) *proto.CheckPoint {
	var cp = &proto.CheckPoint{}
	key := buildCheckPointKey(prefix)
	err := l.db.Get(key, cp)
	if err != nil {
		logger.Log("master", "getCheckPoint", "missing-key", key)
		if l.basePrefix != prefix {
			if opts.Mode == ModeLatest {
				latest := l.getCheckPoint(l.basePrefix, l.opts)
				l.commit(prefix, latest.Index)
			} else if opts.Mode == ModeEarliest {
				l.commit(prefix, 0)
			} else if opts.Mode == ModeCustom {
				max := l.getCheckPoint(l.basePrefix, opts).Index
				if opts.CustomIndex >= max {
					l.commit(prefix, max)
				} else {
					l.commit(prefix, opts.CustomIndex)
				}
			}
		}
	}
	return cp
}

func (l *MasterLedger) commit(prefix string, index uint64) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	cp := &proto.CheckPoint{
		Index: index,
	}
	key := buildCheckPointKey(prefix)
	return l.db.Put(key, cp)
}

func (l *Ledger) Open() {
	startIdx := l.master.getCheckPoint(l.basePrefix, l.opts).Index
	scanPrefix := l.buildWriteScanKey()
	logger.Log("slave", "scanning", "prefix", scanPrefix, "startIdx", startIdx)
	l.master.db.ScanKeysIndexed(scanPrefix, startIdx, func(k []byte, idx uint64) {
		logger.Log("slave", "replaying", "key", k)
		value, err := l.master.db.GetBytes(k)
		if err == nil {
			_, err = l.out.Write(value)
			l.commit(idx)
		} else {
			panic(err)
		}
	})
}

func (l *Ledger) Commit() {
	index := l.master.getCheckPoint(l.master.basePrefix, l.opts).Index
	l.commit(index)
}

func (l *Ledger) commit(index uint64) error {
	cp := &proto.CheckPoint{
		Index: index,
	}
	key := buildCheckPointKey(l.basePrefix)
	return l.master.db.Put(key, cp)
}

func (l *Ledger) buildWriteScanKey() []byte {
	return []byte(fmt.Sprintf("%s-write", l.master.basePrefix))
}

func buildCheckPointKey(prefix string) []byte {
	return []byte(fmt.Sprintf("%s-checkpoint", prefix))
}

func buildWriteSeqKey(prefix string) []byte {
	return []byte(fmt.Sprintf("seq-%s-write", prefix))
}

func buildWriteKey(prefix string, seq uint64) []byte {
	sfmt := "%s-write-%0." + strconv.Itoa(uintDigits) + "d"
	return []byte(fmt.Sprintf(sfmt, prefix, seq))
}
