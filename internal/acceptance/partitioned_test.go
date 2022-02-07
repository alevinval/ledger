package acceptance

import (
	"testing"
	"time"

	"github.com/alevinval/ledger/internal/testutils"
	"github.com/alevinval/ledger/pkg/ledger"
	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/assert"
)

type LedgerWriter interface {
	Write([]byte) (uint64, error)
}

func TestPartitionedWriteAndRead(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		pw, err := ledger.NewPartitionedWriter("channel-1", db, 3)
		assert.Nil(t, err)

		pr, err := pw.NewReader("client-1")
		assert.Nil(t, err)
		defer pr.Close()

		doWrites(t, pw, "first", "second", "third", "fourth")
		assertReadsPartitioned(t, pr, "first", "second", "third", "fourth", "")
	})
}

func TestPartitionedWriteAndReadTwoSequence(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		writer, err := ledger.NewPartitionedWriter("channel-1", db, 3)
		assert.Nil(t, err)

		reader, err := writer.NewReader("client-1")
		assert.Nil(t, err)

		doWrites(t, writer, "first", "second", "third", "fourth")
		assertReadsPartitioned(t, reader, "first", "second", "third", "fourth", "")
		reader.Close()

		reader, err = writer.NewReader("client-1")
		assert.Nil(t, err)

		doWrites(t, writer, "fifth")
		assertReadsPartitioned(t, reader, "fifth", "")
		reader.Close()
	})
}

func doWrites(t *testing.T, w LedgerWriter, data ...string) {
	for _, payload := range data {
		t.Logf("writing %q\n", payload)
		w.Write([]byte(payload))
	}
	time.Sleep(100 * time.Millisecond)
}
