package acceptance

import (
	"testing"
	"time"

	"github.com/alevinval/ledger"
	"github.com/alevinval/ledger/pkg/testutils"
	"github.com/dgraph-io/badger/v2"
	"github.com/stretchr/testify/assert"
)

func TestPartitionedWriteAndRead(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		pw, err := ledger.NewPartitionedWriter("channel-1", db, 3)
		assert.Nil(t, err)

		pr, err := pw.NewReader("client-1")
		assert.Nil(t, err)

		pw.Write([]byte("first message"))
		pw.Write([]byte("second message"))
		pw.Write([]byte("third message"))
		pw.Write([]byte("fourth message"))

		time.Sleep(1 * time.Second)

		assertReadsPartitioned(t, pr, "first message", "second message", "third message", "fourth message", "")
	})
}
