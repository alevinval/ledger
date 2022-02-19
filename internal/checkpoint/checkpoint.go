package checkpoint

import (
	"errors"
	"fmt"
	"sync"

	"github.com/alevinval/ledger/internal/base"
	"github.com/alevinval/ledger/internal/log"
	"github.com/alevinval/ledger/internal/storage"
	"github.com/alevinval/ledger/pkg/proto"
	"go.uber.org/zap"
)

var (
	ErrUnknownOffsetMode = errors.New("unknown offset mode")

	logger = log.GetLogger()
)

type Checkpoint struct {
	storage *storage.Storage
	opts    *base.Options

	mu struct {
		sync.RWMutex
		commitCp *proto.Checkpoint // Instance for zero-allocation commits
	}

	key []byte
}

func NewCheckpoint(basePrefix string, s *storage.Storage, opts *base.Options) *Checkpoint {
	return &Checkpoint{
		key:     buildCheckpointKey(basePrefix),
		storage: s,
		opts:    opts,

		mu: struct {
			sync.RWMutex
			commitCp *proto.Checkpoint
		}{
			commitCp: &proto.Checkpoint{},
		},
	}
}

func (c *Checkpoint) GetCheckpoint() (*proto.Checkpoint, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cp := &proto.Checkpoint{}
	return cp, c.storage.Get(c.key, cp)
}

func (c *Checkpoint) GetCheckpointFrom(other *Checkpoint) (*proto.Checkpoint, error) {
	cp, err := other.GetCheckpoint()
	if err != nil {
		logger.Error("error getting checkpoint from other", zap.ByteString("from", other.key), zap.ByteString("to", c.key), zap.String("error", err.Error()))
		return cp, err
	}

	switch c.opts.Offset {
	case base.LatestOffset:
	case base.EarliestOffset:
		cp.Offset = 0
	case base.CustomOffset:
		if c.opts.CustomOffset < cp.Offset {
			cp.Offset = c.opts.CustomOffset
		}
	default:
		logger.Error("unknown checkpoint mode used", zap.Reflect("value", c.opts.Offset))
		return cp, ErrUnknownOffsetMode
	}
	return cp, err
}

func (c *Checkpoint) Commit(offset uint64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if offset < c.mu.commitCp.Offset {
		logger.Warn("commiting checkpoint backwards", zap.Uint64("from", c.mu.commitCp.Offset), zap.Uint64("to", offset))
	}

	c.mu.commitCp.Offset = offset
	return c.storage.Put(c.key, c.mu.commitCp)
}

func buildCheckpointKey(prefix string) []byte {
	return []byte(fmt.Sprintf("%s-checkpoint", prefix))
}
