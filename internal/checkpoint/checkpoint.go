package checkpoint

import (
	"fmt"
	"sync"

	"github.com/alevinval/ledger/internal/base"
	"github.com/alevinval/ledger/internal/storage"
	"github.com/alevinval/ledger/internal/util/log"
	"github.com/alevinval/ledger/pkg/proto"
	"go.uber.org/zap"
)

var (
	logger = log.GetLogger()
)

type Checkpoint struct {
	key     []byte
	storage *storage.Storage
	opts    *base.Options

	mu struct {
		sync.RWMutex
		commitCp *proto.Checkpoint // Cache instance for comitting without allocation
	}
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
	err := c.storage.Get(c.key, cp)
	return cp, err
}

func (c *Checkpoint) GetCheckpointFrom(other *Checkpoint) (*proto.Checkpoint, error) {
	cp, err := other.GetCheckpoint()
	if err != nil {
		logger.Error("cannot get checkpoint from ", zap.ByteString("from", other.key), zap.ByteString("to", c.key))
		return cp, err
	}

	switch c.opts.Offset {
	case base.LatestOffset:
		// Nothing to do
	case base.EarliestOffset:
		cp.Offset = 0
	case base.CustomOffset:
		if c.opts.CustomOffset < cp.Offset {
			cp.Offset = c.opts.CustomOffset
		}
	default:
		logger.Fatal("unknown checkpoint mode used", zap.String("mode", c.opts.Offset.String()))
	}
	return cp, err
}

func (c *Checkpoint) Commit(offset uint64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.mu.commitCp.Offset = offset
	return c.storage.Put(c.key, c.mu.commitCp)
}

func buildCheckpointKey(prefix string) []byte {
	return []byte(fmt.Sprintf("%s-checkpoint", prefix))
}
