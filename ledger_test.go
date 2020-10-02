package ledger

import (
	"bytes"
	"log"
	"os"
	"path"
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/stretchr/testify/assert"
)

func TestLedgerWriteAndRead(t *testing.T) {
	options := []*Options{
		{Mode: ModeLatest, SequenceBandwidth: 1000, KeySpaceBatchSize: 10},
		{Mode: ModeLatest, SequenceBandwidth: 1000, KeySpaceBatchSize: 100},
		{Mode: ModeLatest, SequenceBandwidth: 1000, KeySpaceBatchSize: 150},
		{Mode: ModeLatest, SequenceBandwidth: 1000, KeySpaceBatchSize: 1000},
		{Mode: ModeLatest, SequenceBandwidth: 1000, KeySpaceBatchSize: 1250},
		{Mode: ModeLatest, SequenceBandwidth: 1000, KeySpaceBatchSize: 5},
		{Mode: ModeLatest, SequenceBandwidth: 1000, KeySpaceBatchSize: 50},
		{Mode: ModeLatest, SequenceBandwidth: 1000, KeySpaceBatchSize: 347},
		{Mode: ModeLatest, SequenceBandwidth: 1000, KeySpaceBatchSize: 500},
		{Mode: ModeLatest, SequenceBandwidth: 1000, KeySpaceBatchSize: 5000},
	}
	for _, opts := range options {
		t.Logf("running test with options: %v", opts)
		runTest(func(db *badger.DB) {
			master := NewMasterOpts("channel-1", db, opts)
			master.Write([]byte("zero"))
			master.Write([]byte("first"))

			out := new(bytes.Buffer)
			slave := NewOpts(master, "client-1", out, opts)
			master.Write([]byte("second"))
			master.Write([]byte("third"))
			slave.Open()
			assert.Equal(t, "secondthird", out.String())

			out = new(bytes.Buffer)
			slave = NewOpts(master, "client-1", out, opts)
			slave.Open()
			assert.Equal(t, "", out.String())
		})

	}
}

func TestLedgerModeEarliest(t *testing.T) {
	runTest(func(db *badger.DB) {
		master := NewMaster("channel-1", db)
		master.Write([]byte("first"))
		master.Write([]byte("second"))

		out := new(bytes.Buffer)
		slave := NewOpts(master, "client-1", out, &Options{Mode: ModeEarliest})
		slave.Open()

		assert.Equal(t, "firstsecond", out.String())
	})
}

func TestLedgerModeCustom(t *testing.T) {
	runTest(func(db *badger.DB) {
		master := NewMaster("channel-1", db)
		master.Write([]byte("first"))
		master.Write([]byte("second"))

		out := new(bytes.Buffer)
		opts := &Options{
			Mode:              ModeCustom,
			CustomIndex:       1,
			KeySpaceBatchSize: 1000,
		}
		slave := NewOpts(master, "client-1", out, opts)
		slave.Open()

		assert.Equal(t, "second", out.String())
	})
}

func TestLedgerModeCustomGreaterThanMax(t *testing.T) {
	runTest(func(db *badger.DB) {
		master := NewMaster("channel-1", db)
		master.Write([]byte("first"))
		master.Write([]byte("second"))

		out := new(bytes.Buffer)
		opts := &Options{
			Mode:              ModeCustom,
			CustomIndex:       10,
			KeySpaceBatchSize: 1000,
		}
		slave := NewOpts(master, "client-1", out, opts)

		master.Write([]byte("third"))
		slave.Open()

		assert.Equal(t, "third", out.String())
	})
}

func openBadgerDB() (*badger.DB, error) {
	tmpPath := os.TempDir()
	storePath := path.Join(tmpPath, "test-badger.db")

	os.RemoveAll(storePath)

	opts := badger.DefaultOptions("")
	opts.Dir = storePath
	opts.ValueDir = storePath
	opts.Logger = nil

	log.Printf("opening badger db in %s", storePath)
	return badger.Open(opts)
}

func runTest(fn func(db *badger.DB)) {
	db, err := openBadgerDB()
	defer db.Close()
	if err != nil {
		log.Fatalf("cannot open badger db: %s", err)
	}

	fn(db)
}
