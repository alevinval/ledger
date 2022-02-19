package acceptance

import (
	"testing"

	"github.com/alevinval/ledger/internal/testutils"
	"github.com/alevinval/ledger/pkg/ledger"
	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/assert"
)

func TestPartitionedWriteAndRead(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		pw, err := ledger.NewPartitionedWriter("channel-1", db, 3)
		assert.Nil(t, err)

		pr, err := pw.NewReader("client-1")
		assert.Nil(t, err)
		defer pr.Close()

		testutils.AssertWrites(t, pw, "first", "second", "third", "fourth")
		testutils.AssertReads(t, pr, "first", "second", "third", "fourth", "")
	})
}

func TestPartitionedWriteAndReadTwoSequence(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		writer, err := ledger.NewPartitionedWriter("channel-1", db, 3)
		assert.Nil(t, err)

		reader, err := writer.NewReader("client-1")
		assert.Nil(t, err)

		testutils.AssertWrites(t, writer, "first", "second", "third", "fourth")
		testutils.AssertReads(t, reader, "first", "second", "third", "fourth", "")
		reader.Close()

		reader, err = writer.NewReader("client-1")
		assert.Nil(t, err)

		testutils.AssertWrites(t, writer, "fifth")
		testutils.AssertReads(t, reader, "fifth", "")
		reader.Close()
	})
}
