package ledger

import (
	"testing"

	"github.com/alevinval/ledger/pkg/proto"
	"github.com/alevinval/ledger/pkg/testutils"
	"github.com/dgraph-io/badger/v2"
	"github.com/stretchr/testify/assert"
)

func TestCheckpointKey(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		opts := DefaultOptions()
		s := &storage{db, opts}

		chk := newCheckpoint("prefix", s, opts)

		assert.Equal(t, "prefix-checkpoint", string(chk.key))
	})
}

func TestCheckpointGetMissing(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		opts := DefaultOptions()
		s := &storage{db, opts}

		chk := newCheckpoint("prefix", s, opts)
		cp, err := chk.GetCheckpoint()

		assert.Equal(t, &proto.Checkpoint{}, cp)
		assert.Error(t, err)
	})
}

func TestCheckpointCommit(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		opts := DefaultOptions()
		s := &storage{db, opts}

		chk := newCheckpoint("prefix", s, opts)
		err := chk.Commit(123)

		assert.NoError(t, err)
	})
}

func TestCheckpointCommitAndGet(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		opts := DefaultOptions()
		s := &storage{db, opts}

		chk := newCheckpoint("prefix", s, opts)
		chk.Commit(123)
		cp, err := chk.GetCheckpoint()

		assert.NoError(t, err)
		assert.Equal(t, uint64(123), cp.Offset)
	})
}

func TestCheckpointFromEarliest(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		opts := DefaultOptions()
		opts.Offset = EarliestOffset
		s := &storage{db, opts}

		master := newCheckpoint("prefix-master", s, opts)
		master.Commit(10)
		reader := newCheckpoint("prefix-reader", s, opts)

		cp, err := reader.GetCheckpointFrom(master)
		assert.NoError(t, err)
		assert.Equal(t, uint64(0), cp.Offset)
	})
}

func TestCheckpointFromLatest(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		opts := DefaultOptions()
		opts.Offset = LatestOffset
		s := &storage{db, opts}

		master := newCheckpoint("prefix-master", s, opts)
		master.Commit(10)
		reader := newCheckpoint("prefix-reader", s, opts)

		cp, err := reader.GetCheckpointFrom(master)
		assert.NoError(t, err)
		assert.Equal(t, uint64(10), cp.Offset)
	})
}

func TestCheckpointFromCustom(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		opts := DefaultOptions()
		opts.Offset = CustomOffset
		opts.CustomOffset = 5
		s := &storage{db, opts}

		master := newCheckpoint("prefix-master", s, opts)
		master.Commit(10)
		reader := newCheckpoint("prefix-reader", s, opts)

		cp, err := reader.GetCheckpointFrom(master)
		assert.NoError(t, err)
		assert.Equal(t, uint64(5), cp.Offset)
	})
}
