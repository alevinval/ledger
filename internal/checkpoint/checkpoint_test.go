package checkpoint

import (
	"testing"

	"github.com/alevinval/ledger/internal/base"
	"github.com/alevinval/ledger/internal/storage"
	"github.com/alevinval/ledger/internal/testutils"
	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

var defaultOpts = base.DefaultOptions()

func captureLogs(level zapcore.Level) (logs *observer.ObservedLogs, restore func()) {
	originalLogger := logger
	restoreOriginalLogger := func() {
		logger = originalLogger
	}

	core, logs := observer.New(zap.WarnLevel)
	newLogger := zap.New(core)
	logger = newLogger
	return logs, restoreOriginalLogger
}

func TestCheckpoint_keyMatchesPrefix(t *testing.T) {
	withCheckpoint("prefix", defaultOpts, func(checkpoint *Checkpoint) {
		assert.Equal(t, "prefix-checkpoint", string(checkpoint.key))
	})
}

func TestCheckpointGetCheckpoint_missingInStore(t *testing.T) {
	withCheckpoint("prefix", defaultOpts, func(checkpoint *Checkpoint) {
		_, err := checkpoint.GetCheckpoint()

		assert.ErrorIs(t, err, badger.ErrKeyNotFound)
	})
}

func TestCheckpointGetCheckpointFrom_errorOtherNotCommitted(t *testing.T) {
	withWriterAndReader(defaultOpts, func(writer, reader *Checkpoint) {
		_, err := reader.GetCheckpointFrom(writer)

		assert.ErrorIs(t, err, badger.ErrKeyNotFound)
	})
}

func TestCheckpointCommit(t *testing.T) {
	withCheckpoint("prefix", defaultOpts, func(checkpoint *Checkpoint) {
		err := checkpoint.Commit(123)

		assert.NoError(t, err)
	})
}

func TestCheckpointCommit_getCheckpointGivesCommittedValue(t *testing.T) {
	withCheckpoint("prefix", defaultOpts, func(checkpoint *Checkpoint) {
		checkpoint.Commit(123)
		cp, err := checkpoint.GetCheckpoint()

		assert.NoError(t, err)
		assert.Equal(t, uint64(123), cp.Offset)
	})
}

func TestCheckpointCommit_warnBackwardsCommit(t *testing.T) {
	logs, restore := captureLogs(zap.WarnLevel)
	defer restore()

	withCheckpoint("prefix", defaultOpts, func(checkpoint *Checkpoint) {
		checkpoint.Commit(123)
		checkpoint.Commit(110)

		assert.Greater(t, logs.Len(), 0)

		backwardsWarn := logs.All()[0]
		assert.Equal(t, "commiting checkpoint backwards", backwardsWarn.Message)
		assert.Equal(t, uint64(123), backwardsWarn.ContextMap()["from"])
		assert.Equal(t, uint64(110), backwardsWarn.ContextMap()["to"])
	})
}

func TestCheckpointGetCheckpointFrom_errorUnknownOffset(t *testing.T) {
	opts := base.DefaultOptions()
	opts.Offset = (base.OffsetMode(99))

	withWriterAndReader(opts, func(writer, reader *Checkpoint) {
		writer.Commit(1)
		_, err := reader.GetCheckpointFrom(writer)

		assert.ErrorIs(t, err, ErrUnknownOffsetMode)
	})
}

func TestCheckpointGetCheckpointFrom_usingEarliestOffset(t *testing.T) {
	opts := base.DefaultOptions()
	opts.Offset = base.EarliestOffset

	withWriterAndReader(opts, func(writer, reader *Checkpoint) {
		writer.Commit(10)
		cp, err := reader.GetCheckpointFrom(writer)

		assert.NoError(t, err)
		assert.Equal(t, uint64(0), cp.Offset)
	})
}

func TestCheckpointGetCheckpointFrom_usingLatestOffset(t *testing.T) {
	opts := base.DefaultOptions()
	opts.Offset = base.LatestOffset

	withWriterAndReader(opts, func(writer, reader *Checkpoint) {
		writer.Commit(10)
		cp, err := reader.GetCheckpointFrom(writer)

		assert.NoError(t, err)
		assert.Equal(t, uint64(10), cp.Offset)
	})
}

func TestCheckpointGetCheckpointFrom_usingCustomOffset(t *testing.T) {
	opts := base.DefaultOptions()
	opts.Offset = base.CustomOffset
	opts.CustomOffset = 5

	withWriterAndReader(opts, func(writer, reader *Checkpoint) {
		writer.Commit(10)
		cp, err := reader.GetCheckpointFrom(writer)

		assert.NoError(t, err)
		assert.Equal(t, uint64(5), cp.Offset)
	})
}

func withWriterAndReader(opts *base.Options, fn func(writer *Checkpoint, reader *Checkpoint)) {
	withCheckpoint("writer", opts, func(writer *Checkpoint) {
		withCheckpoint("reader", opts, func(reader *Checkpoint) {
			fn(writer, reader)
		})
	})
}

func withCheckpoint(basePrefix string, opts *base.Options, fn func(c *Checkpoint)) {
	withStorage(opts, func(storage *storage.Storage) {
		checkpoint := NewCheckpoint(basePrefix, storage, opts)
		fn(checkpoint)
	})
}

func withStorage(opts *base.Options, fn func(storage *storage.Storage)) {
	testutils.WithDB(func(db *badger.DB) {
		storage := &storage.Storage{DB: db, Opts: opts}
		fn(storage)
	})
}
