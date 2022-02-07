package checkpoint

import (
	"testing"

	"github.com/alevinval/ledger/internal/base"
	"github.com/alevinval/ledger/internal/storage"
	"github.com/alevinval/ledger/internal/testutils"
	"github.com/alevinval/ledger/pkg/proto"
	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/assert"
)

func WithStorage(opts *base.Options, fn func(storage *storage.Storage)) {
	testutils.WithDB(func(db *badger.DB) {
		storage := &storage.Storage{DB: db, Opts: opts}
		fn(storage)
	})
}

func TestCheckpointKey(t *testing.T) {
	opts := base.DefaultOptions()

	WithStorage(opts, func(storage *storage.Storage) {
		chk := NewCheckpoint("prefix", storage, opts)

		assert.Equal(t, "prefix-checkpoint", string(chk.key))
	})
}

func TestCheckpointGetMissing(t *testing.T) {
	opts := base.DefaultOptions()

	WithStorage(opts, func(storage *storage.Storage) {
		chk := NewCheckpoint("prefix", storage, opts)
		cp, err := chk.GetCheckpoint()

		assert.Equal(t, &proto.Checkpoint{}, cp)
		assert.Error(t, err)
	})
}

func TestCheckpointCommit(t *testing.T) {
	opts := base.DefaultOptions()

	WithStorage(opts, func(storage *storage.Storage) {
		chk := NewCheckpoint("prefix", storage, opts)
		err := chk.Commit(123)

		assert.NoError(t, err)
	})
}

func TestCheckpointCommitAndGet(t *testing.T) {
	opts := base.DefaultOptions()

	WithStorage(opts, func(storage *storage.Storage) {
		chk := NewCheckpoint("prefix", storage, opts)
		chk.Commit(123)
		cp, err := chk.GetCheckpoint()

		assert.NoError(t, err)
		assert.Equal(t, uint64(123), cp.Offset)
	})
}

func TestCheckpointFromEarliest(t *testing.T) {
	opts := base.DefaultOptions()
	opts.Offset = base.EarliestOffset

	WithStorage(opts, func(storage *storage.Storage) {
		master := NewCheckpoint("prefix-master", storage, opts)
		master.Commit(10)
		reader := NewCheckpoint("prefix-reader", storage, opts)

		cp, err := reader.GetCheckpointFrom(master)
		assert.NoError(t, err)
		assert.Equal(t, uint64(0), cp.Offset)
	})
}

func TestCheckpointFromLatest(t *testing.T) {
	opts := base.DefaultOptions()
	opts.Offset = base.LatestOffset

	WithStorage(opts, func(storage *storage.Storage) {
		master := NewCheckpoint("prefix-master", storage, opts)
		master.Commit(10)
		reader := NewCheckpoint("prefix-reader", storage, opts)

		cp, err := reader.GetCheckpointFrom(master)
		assert.NoError(t, err)
		assert.Equal(t, uint64(10), cp.Offset)
	})
}

func TestCheckpointFromCustom(t *testing.T) {
	opts := base.DefaultOptions()
	opts.Offset = base.CustomOffset
	opts.CustomOffset = 5

	WithStorage(opts, func(storage *storage.Storage) {
		master := NewCheckpoint("prefix-master", storage, opts)
		master.Commit(10)
		reader := NewCheckpoint("prefix-reader", storage, opts)

		cp, err := reader.GetCheckpointFrom(master)
		assert.NoError(t, err)
		assert.Equal(t, uint64(5), cp.Offset)
	})
}
