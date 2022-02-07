package checkpoint

import (
	"fmt"

	"github.com/alevinval/ledger/internal/base"
	"github.com/alevinval/ledger/internal/storage"
	"github.com/alevinval/ledger/internal/util/log"
	"github.com/alevinval/ledger/pkg/proto"
	"go.uber.org/zap"
)

var (
	logger = log.GetLogger()
)

const (
	// CustomOffset reads from a specific offset onwards.
	CustomOffset OffsetMode = iota
	// EarliestOffset reads from earliest possible offset onwards.
	EarliestOffset
)

type OffsetMode byte

type Checkpoint struct {
	key     []byte
	storage *storage.Storage
	opts    *base.Options
}

func NewCheckpoint(basePrefix string, s *storage.Storage, opts *base.Options) *Checkpoint {
	return &Checkpoint{
		key:     buildCheckpointKey(basePrefix),
		storage: s,
		opts:    opts,
	}
}

func (c *Checkpoint) GetCheckpoint() (*proto.Checkpoint, error) {
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
	case base.EarliestOffset:
		cp.Offset = 0
	case base.CustomOffset:
		if c.opts.CustomOffset < cp.Offset {
			cp.Offset = c.opts.CustomOffset
		}
	case base.LatestOffset:
		// Nothing to do
	default:
		logger.Fatal("unknown checkpoint mode used", zap.String("mode", c.opts.Offset.String()))
	}
	return cp, err
}

func (c *Checkpoint) Commit(offset uint64) error {
	cp := &proto.Checkpoint{
		Offset: offset,
	}
	return c.storage.Put(c.key, cp)
}

func buildCheckpointKey(prefix string) []byte {
	return []byte(fmt.Sprintf("%s-checkpoint", prefix))
}
