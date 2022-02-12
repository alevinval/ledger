package ledger

import (
	"testing"

	"github.com/alevinval/ledger/internal/testutils"
	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/assert"
)

func TestWriterWrite(t *testing.T) {
	withWriter(t, func(w *Writer) {
		_, err := w.Write([]byte("some value"))

		assert.NoError(t, err)
	})
}

func TestWriterWrite_returnsWriteOffset(t *testing.T) {
	withWriter(t, func(w *Writer) {
		offset, err := w.Write([]byte("some value"))

		assert.NoError(t, err)
		assert.Equal(t, uint64(1), offset)
	})
}

func TestWriterWrite_advancesOffset(t *testing.T) {
	withWriter(t, func(w *Writer) {
		testutils.AssertCheckpointAt(t, w, 0)

		w.Write([]byte("some value"))
		w.Write([]byte("some value"))
		lastOffset, _ := w.Write([]byte("some value"))

		testutils.AssertCheckpointAt(t, w, 3)
		assert.Equal(t, uint64(3), lastOffset)
	})
}

func TestWriterWrite_failsOnClosedWriter(t *testing.T) {
	withWriter(t, func(w *Writer) {
		w.Close()

		_, err := w.Write([]byte("some-value"))
		assert.ErrorIs(t, err, ErrClosedWriter)
		testutils.AssertCheckpointAt(t, w, 0)
	})
}

func TestNewWriter_failsWithClosedDB(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		db.Close()

		_, err := NewWriter("writer", db)
		assert.ErrorIs(t, err, badger.ErrDBClosed)
	})
}

func TestWriterWrite_failsWithClosedDB(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		writer, err := NewWriter("writer", db)
		assert.NoError(t, err)

		db.Close()

		_, err = writer.Write([]byte("some-value"))
		assert.ErrorIs(t, err, badger.ErrDBClosed)
	})
}

func withWriter(t *testing.T, fn func(w *Writer)) {
	testutils.WithDB(func(db *badger.DB) {
		writer, err := NewWriter("writer", db)
		assert.NoError(t, err)

		defer writer.Close()

		fn(writer)
	})
}
