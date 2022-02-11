package acceptance

import (
	"testing"

	"github.com/alevinval/ledger/internal/testutils"
	"github.com/alevinval/ledger/pkg/ledger"
	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/assert"
)

func TestLedgerLatestOffset(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		opts := ledger.DefaultOptions()
		w, err := ledger.NewWriterOpts("channel-1", db, opts)
		assert.Nil(t, err)
		defer w.Close()

		w.Write([]byte("zero"))
		w.Write([]byte("first"))

		r, err := w.NewReader("client-1")
		assert.Nil(t, err)
		defer r.Close()

		w.Write([]byte("second"))
		w.Write([]byte("third"))

		testutils.AssertReads(t, r, "second", "third", "")
	})
}

func TestLedgerEarliestOffset(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		w, err := ledger.NewWriter("channel-1", db)
		assert.Nil(t, err)
		defer w.Close()

		w.Write([]byte("first"))
		w.Write([]byte("second"))

		opts := ledger.DefaultOptions()
		opts.Offset = ledger.EarliestOffset
		r, err := w.NewReaderOpts("client-1", opts)
		assert.Nil(t, err)
		defer r.Close()
		w.Write([]byte("third"))
		testutils.AssertReads(t, r, "first", "second", "third", "")
	})
}

func TestLedgerCustomOffset(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		w, err := ledger.NewWriter("channel-1", db)
		assert.Nil(t, err)
		defer w.Close()

		w.Write([]byte("first"))
		w.Write([]byte("second"))

		opts := ledger.DefaultOptions()
		opts.Offset = ledger.CustomOffset
		opts.CustomOffset = 1
		r, err := w.NewReaderOpts("client-1", opts)
		assert.Nil(t, err)
		defer r.Close()
		testutils.AssertReads(t, r, "second", "")
	})
}

func TestLedgerCustomOffsetGreaterThanMax(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		w, err := ledger.NewWriter("channel-1", db)
		assert.Nil(t, err)
		defer w.Close()

		w.Write([]byte("first"))
		w.Write([]byte("second"))

		opts := ledger.DefaultOptions()
		opts.Offset = ledger.CustomOffset
		opts.CustomOffset = 10
		r, err := w.NewReaderOpts("client-1", opts)
		assert.Nil(t, err)
		defer r.Close()

		w.Write([]byte("third"))

		testutils.AssertReads(t, r, "third", "")
	})
}

func TestLedgerReadWriteReadWithoutAutoCommit(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		w, err := ledger.NewWriter("channel-1", db)
		assert.Nil(t, err)

		r, err := w.NewReader("client-1")
		assert.Nil(t, err)

		w.Write([]byte("first"))
		w.Write([]byte("second"))

		testutils.AssertReadsNoCommit(t, r, "first", "second", "")

		w.Write([]byte("third"))
		w.Write([]byte("fourth"))

		/// Even if there is no commit, it resumes from last received message
		testutils.AssertReadsNoCommit(t, r, "third", "fourth", "")

		// Because the previous reads happened without a commit,
		// a new reader starts from last known commit
		r, err = w.NewReader("client-1")

		testutils.AssertReadsNoCommit(t, r, "first", "second", "third", "fourth", "")
	})
}

func TestLedgerContinuesFromLastCheckpoint(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		w, err := ledger.NewWriter("channel-1", db)
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

		testutils.AssertReads(t, r, "first", "")
	})
}

func TestLedgerReadGreaterThanBatchSize(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		opts := ledger.DefaultOptions()
		opts.BatchSize = 10
		w, err := ledger.NewWriterOpts("channel-1", db, opts)
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
			total += string(s.Data())
		}

		assert.Equal(t, int(3*opts.BatchSize), len(total))
	})
}

func TestLedgerClose(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		w, err := ledger.NewWriter("writer-1", db)
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
