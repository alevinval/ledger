package ledger

import (
	"bytes"
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
			w, err := NewWriterOpts("channel-1", db, opts)
			assert.Nil(t, err)
			w.Write([]byte("zero"))
			w.Write([]byte("first"))

			out := new(bytes.Buffer)
			r, err := NewReaderOpts(w, "client-1", out, opts)
			assert.Nil(t, err)
			w.Write([]byte("second"))
			w.Write([]byte("third"))
			r.Open()
			assert.Equal(t, "secondthird", out.String())

			out = new(bytes.Buffer)
			r, err = NewReaderOpts(w, "client-1", out, opts)
			assert.Nil(t, err)
			r.Open()
			assert.Equal(t, "", out.String())
		})

	}
}

func TestLedgerModeEarliest(t *testing.T) {
	runTest(func(db *badger.DB) {
		w, err := NewWriter("channel-1", db)
		assert.Nil(t, err)
		w.Write([]byte("first"))
		w.Write([]byte("second"))

		out := new(bytes.Buffer)
		opts := DefaultOptions()
		opts.Mode = ModeEarliest
		r, err := NewReaderOpts(w, "client-1", out, opts)
		assert.Nil(t, err)
		w.Write([]byte("third"))
		r.Open()

		assert.Equal(t, "firstsecondthird", out.String())
	})
}

func TestLedgerModeCustom(t *testing.T) {
	runTest(func(db *badger.DB) {
		w, err := NewWriter("channel-1", db)
		assert.Nil(t, err)
		w.Write([]byte("first"))
		w.Write([]byte("second"))

		out := new(bytes.Buffer)
		opts := DefaultOptions()
		opts.Mode = ModeCustom
		opts.CustomIndex = 1
		r, err := NewReaderOpts(w, "client-1", out, opts)
		assert.Nil(t, err)
		r.Open()

		assert.Equal(t, "second", out.String())
	})
}

func TestLedgerModeCustomGreaterThanMax(t *testing.T) {
	runTest(func(db *badger.DB) {
		w, err := NewWriter("channel-1", db)
		assert.Nil(t, err)
		w.Write([]byte("first"))
		w.Write([]byte("second"))

		out := new(bytes.Buffer)
		opts := DefaultOptions()
		opts.Mode = ModeCustom
		opts.CustomIndex = 10
		r, err := NewReaderOpts(w, "client-1", out, opts)
		assert.Nil(t, err)

		w.Write([]byte("third"))
		r.Open()

		assert.Equal(t, "third", out.String())
	})
}

func TestLedgerOpenTicker(t *testing.T) {
	runTest(func(db *badger.DB) {
		w, err := NewWriter("channel-1", db)
		assert.Nil(t, err)

		out := new(bytes.Buffer)
		r, err := NewReader(w, "client-1", out)
		assert.Nil(t, err)

		w.Write([]byte("first"))
		r.OpenTicker(400)
		assert.Equal(t, "first", out.String())

		time.Sleep(300 * time.Millisecond)
		w.Write([]byte("second"))
		assert.Equal(t, "first", out.String())

		time.Sleep(150 * time.Millisecond)
		assert.Equal(t, "firstsecond", out.String())
	})
}

func TestLedgerMoreThanOneBatchSize(t *testing.T) {
	runTest(func(db *badger.DB) {
		w, err := NewWriter("channel-1", db)
		assert.Nil(t, err)

		out := new(bytes.Buffer)
		r, err := NewReader(w, "client-1", out)
		assert.Nil(t, err)

		r.OpenTicker(100)

		opts := DefaultOptions()
		for i := 0; i < int(3*opts.KeySpaceBatchSize); i++ {
			w.Write([]byte("1"))
		}
		time.Sleep(200 * time.Millisecond)

		assert.Equal(t, int(3*opts.KeySpaceBatchSize), len(out.Bytes()))
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
