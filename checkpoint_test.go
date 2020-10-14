package ledger

import (
	"testing"

	"github.com/alevinval/ledger/pkg/proto"
	"github.com/dgraph-io/badger/v2"
	"github.com/stretchr/testify/assert"
)

func TestCheckpointKey(t *testing.T) {
	runTest(func(db *badger.DB) {
		opts := DefaultOptions()
		s := &storage{db, opts}

		chk := newCheckpoint("prefix", s, opts)

		assert.Equal(t, "prefix-checkpoint", string(chk.key))
	})
}

func TestCheckpointGetMissing(t *testing.T) {
	runTest(func(db *badger.DB) {
		opts := DefaultOptions()
		s := &storage{db, opts}

		chk := newCheckpoint("prefix", s, opts)
		cp, err := chk.GetCheckpoint()

		assert.Equal(t, &proto.Checkpoint{}, cp)
		assert.Error(t, err)
	})
}

func TestCheckpointCommit(t *testing.T) {
	runTest(func(db *badger.DB) {
		opts := DefaultOptions()
		s := &storage{db, opts}

		chk := newCheckpoint("prefix", s, opts)
		err := chk.Commit(123)

		assert.NoError(t, err)
	})
}

func TestCheckpointCommitAndGet(t *testing.T) {
	runTest(func(db *badger.DB) {
		opts := DefaultOptions()
		s := &storage{db, opts}

		chk := newCheckpoint("prefix", s, opts)
		chk.Commit(123)
		cp, err := chk.GetCheckpoint()

		assert.NoError(t, err)
		assert.Equal(t, uint64(123), cp.Offset)
	})
}
