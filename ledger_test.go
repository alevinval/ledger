package ledger

import (
	"log"
	"os"
	"path"
	"testing"
	"time"

	"github.com/dgraph-io/badger/v2"
	"github.com/stretchr/testify/assert"
)

func TestLedgerWriteAndReadDifferentBatchSizes(t *testing.T) {
	options := []*Options{
		{BatchSize: 10},
		{BatchSize: 100},
		{BatchSize: 150},
		{BatchSize: 1000},
		{BatchSize: 1250},
		{BatchSize: 5},
		{BatchSize: 50},
		{BatchSize: 347},
		{BatchSize: 500},
		{BatchSize: 5000},
	}
	for _, opts := range options {
		opts.SequenceBandwidth = 1000
		opts.Offset = LatestOffset

		t.Logf("running test with options: %v", opts)
		runTest(func(db *badger.DB) {
			w, err := NewWriterOpts("channel-1", db, opts)
			assert.Nil(t, err)
			defer w.Close()

			w.Write([]byte("zero"))
			w.Write([]byte("first"))

			r, err := w.NewReader("client-1")
			assert.Nil(t, err)
			defer r.Close()

			w.Write([]byte("second"))
			w.Write([]byte("third"))

			assertReads(t, r, "second", "third", "")

			r, err = w.NewReader("client-1")
			assert.Nil(t, err)
			defer r.Close()
			assertReads(t, r, "")
		})
	}
}

func TestLedgerRenewInstanceWithNoAutoCommit(t *testing.T) {
	runTest(func(db *badger.DB) {
		w, err := NewWriter("channel-1", db)
		assert.Nil(t, err)

		r, err := w.NewReader("client-1")
		assert.Nil(t, err)

		w.Write([]byte("first"))
		w.Write([]byte("second"))

		assertReadsNoCommit(t, r, "first", "second", "")

		w.Write([]byte("third"))
		w.Write([]byte("fourth"))

		/// Even if there is no commit, it resumes from last received message
		assertReadsNoCommit(t, r, "third", "fourth", "")

		// Because the previous reads happened without a commit,
		// a new reader starts from last known commit
		r, err = w.NewReader("client-1")

		assertReadsNoCommit(t, r, "first", "second", "third", "fourth", "")
	})
}

func TestLedgerModeEarliest(t *testing.T) {
	runTest(func(db *badger.DB) {
		w, err := NewWriter("channel-1", db)
		assert.Nil(t, err)
		defer w.Close()

		w.Write([]byte("first"))
		w.Write([]byte("second"))

		opts := DefaultOptions()
		opts.Offset = EarliestOffset
		r, err := w.NewReaderOpts("client-1", opts)
		assert.Nil(t, err)
		defer r.Close()
		w.Write([]byte("third"))
		assertReads(t, r, "first", "second", "third", "")
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
		opts.Offset = CustomOffset
		opts.CustomOffset = 1
		r, err := w.NewReaderOpts("client-1", opts)
		assert.Nil(t, err)
		defer r.Close()
		assertReads(t, r, "second", "")
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
		opts.Offset = CustomOffset
		opts.CustomOffset = 10
		r, err := w.NewReaderOpts("client-1", opts)
		assert.Nil(t, err)
		defer r.Close()

		w.Write([]byte("third"))

		assertReads(t, r, "third", "")
	})
}

func TestLedgerReconnectAndContinueFromWhereLeft(t *testing.T) {
	runTest(func(db *badger.DB) {
		w, err := NewWriter("channel-1", db)
		assert.Nil(t, err)
		defer w.Close()

		r, err := w.NewReader("client-1")
		assert.Nil(t, err)

		// We've now established a checkpoint for the reader.
		// Close it.
		r.Close()

		w.Write([]byte("first"))

		r, err = w.NewReader("client-1")
		assert.Nil(t, err)

		assertReads(t, r, "first", "")
	})
}

func TestLedgerMoreThanOneBatchSize(t *testing.T) {
	runTest(func(db *badger.DB) {
		opts := DefaultOptions()
		opts.BatchSize = 10
		w, err := NewWriterOpts("channel-1", db, opts)
		assert.Nil(t, err)
		defer w.Close()

		r, err := w.NewReaderOpts("client-1", opts)
		assert.Nil(t, err)
		defer r.Close()

		for i := 0; i < int(3*opts.BatchSize); i++ {
			w.Write([]byte("1"))
		}

		total := ""
		messageCh, err := r.Read()
		assert.Nil(t, err)
		for i := 0; i < int(3*opts.BatchSize); i++ {
			s := <-messageCh
			total += string(s.Data)
		}

		assert.Equal(t, int(3*opts.BatchSize), len(total))
	})
}

func TestLedgerClose(t *testing.T) {
	runTest(func(db *badger.DB) {
		w, err := NewWriter("writer-1", db)
		assert.Nil(t, err)

		r1, err := w.NewReader("client-1")
		assert.Nil(t, err)

		r2, err := w.NewReader("client-2")
		assert.Nil(t, err)

		w.Close()
		r1.Close()
		r2.Close()
	})
}

// assert reads with auto-commits between reads.
func assertReads(t *testing.T, r *Reader, expected ...string) {
	assertReadsImpl(t, r, true, expected...)
}

// assert reads without auto-commits between reads.
func assertReadsNoCommit(t *testing.T, r *Reader, expected ...string) {
	assertReadsImpl(t, r, false, expected...)
}

func assertReadsImpl(t *testing.T, r *Reader, autoCommit bool, expected ...string) {
	ch, err := r.Read()
	assert.Nil(t, err)
	for i := range expected {
		select {
		case actual := <-ch:
			assert.Equal(t, expected[i], string(actual.Data))
			if autoCommit {
				r.Commit(actual.Offset)
			}
		case <-time.After(100 * time.Millisecond):
			if expected[i] != "" {
				assert.FailNowf(t, "missing-read", "expected %q, instead read timeout", expected[i])
			}
		}
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
