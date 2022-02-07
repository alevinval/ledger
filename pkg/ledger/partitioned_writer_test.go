package ledger

import (
	"testing"

	"github.com/alevinval/ledger/internal/testutils"
	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/assert"
)

func TestPartitionedWriter(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		_, err := NewPartitionedWriter("test-partitioned-writer", db, 4)
		assert.Nil(t, err)
	})
}
