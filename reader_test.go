package ledger

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"
	"time"

	"github.com/dgraph-io/badger/v2"
	"github.com/stretchr/testify/assert"
)

func TestLedgerWriteAndRead(t *testing.T) {
	options := []*Options{
		{Mode: ModeLatest, SequenceBandwidth: 1000, BatchSize: 10},
		{Mode: ModeLatest, SequenceBandwidth: 1000, BatchSize: 100},
		{Mode: ModeLatest, SequenceBandwidth: 1000, BatchSize: 150},
		{Mode: ModeLatest, SequenceBandwidth: 1000, BatchSize: 1000},
		{Mode: ModeLatest, SequenceBandwidth: 1000, BatchSize: 1250},
		{Mode: ModeLatest, SequenceBandwidth: 1000, BatchSize: 5},
		{Mode: ModeLatest, SequenceBandwidth: 1000, BatchSize: 50},
		{Mode: ModeLatest, SequenceBandwidth: 1000, BatchSize: 347},
		{Mode: ModeLatest, SequenceBandwidth: 1000, BatchSize: 500},
		{Mode: ModeLatest, SequenceBandwidth: 1000, BatchSize: 5000},
	}
	for _, opts := range options {
		t.Logf("running test with options: %v", opts)
		runTest(func(db *badger.DB) {
			w, err := NewWriterOpts("channel-1", db, opts)
			assert.Nil(t, err)
			defer w.Close()

			w.Write([]byte("zero"))
			w.Write([]byte("first"))

			r, err := NewReaderOpts(w, "client-1", opts)
			assert.Nil(t, err)
			w.Write([]byte("second"))
			w.Write([]byte("third"))
			assertReads(t, r, "second", "third", "")

			r, err = NewReaderOpts(w, "client-1", opts)
			assert.Nil(t, err)
			assertReads(t, r, "")
		})

	}
}

func TestLedgerModeEarliest(t *testing.T) {
	runTest(func(db *badger.DB) {
		w, err := NewWriter("channel-1", db)
		assert.Nil(t, err)
		defer w.Close()

		w.Write([]byte("first"))
		w.Write([]byte("second"))

		opts := DefaultOptions()
		opts.Mode = ModeEarliest
		r, err := NewReaderOpts(w, "client-1", opts)
		assert.Nil(t, err)
		w.Write([]byte("third"))
		assertReads(t, r, "first", "second", "third")
	})
}

func TestLedgerModeCustom(t *testing.T) {
	runTest(func(db *badger.DB) {
		w, err := NewWriter("channel-1", db)
		assert.Nil(t, err)
		defer w.Close()

		w.Write([]byte("first"))
		w.Write([]byte("second"))

		opts := DefaultOptions()
		opts.Mode = ModeCustom
		opts.CustomIndex = 1
		r, err := NewReaderOpts(w, "client-1", opts)
		assert.Nil(t, err)
		assertReads(t, r, "second")
	})
}

func TestLedgerModeCustomGreaterThanMax(t *testing.T) {
	runTest(func(db *badger.DB) {
		w, err := NewWriter("channel-1", db)
		assert.Nil(t, err)
		defer w.Close()

		w.Write([]byte("first"))
		w.Write([]byte("second"))

		opts := DefaultOptions()
		opts.Mode = ModeCustom
		opts.CustomIndex = 10
		r, err := NewReaderOpts(w, "client-1", opts)
		assert.Nil(t, err)

		w.Write([]byte("third"))

		assertReads(t, r, "third")
	})
}

func TestLedgerMoreThanOneBatchSize(t *testing.T) {
	runTest(func(db *badger.DB) {
		w, err := NewWriter("channel-1", db)
		assert.Nil(t, err)
		defer w.Close()

		r, err := NewReader(w, "client-1")
		assert.Nil(t, err)

		opts := DefaultOptions()
		for i := 0; i < int(3*opts.BatchSize); i++ {
			w.Write([]byte("1"))
		}
		time.Sleep(200 * time.Millisecond)

		total := ""
		for i := 0; i < int(3*opts.BatchSize); i++ {
			s, _ := ioutil.ReadAll(r)
			total += string(s)
		}

		s, _ := ioutil.ReadAll(r)
		assert.Equal(t, "", string(s))

		assert.Equal(t, int(3*opts.BatchSize), len(total))
	})
}

func assertReads(t *testing.T, r *Reader, expected ...string) {
	for i := range expected {
		actual, _ := ioutil.ReadAll(r)
		assert.Equal(t, expected[i], string(actual))
	}
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
