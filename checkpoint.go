package ledger

import (
	"fmt"

	"github.com/alevinval/ledger/pkg/proto"
)

type checkpoint struct {
	key  []byte
	db   *storage
	opts *Options
}

func newCheckpoint(basePrefix string, db *storage, opts *Options) *checkpoint {
	return &checkpoint{
		key:  buildCheckpointKey(basePrefix),
		db:   db,
		opts: opts,
	}
}

func (c *checkpoint) GetCheckpoint() (*proto.Checkpoint, error) {
	var cp = &proto.Checkpoint{}
	err := c.db.Get(c.key, cp)
	return cp, err
}

func (c *checkpoint) GetCheckpointFrom(other *checkpoint) (*proto.Checkpoint, error) {
	cp, err := other.GetCheckpoint()
	switch c.opts.Offset {
	case EarliestOffset:
		cp.Offset = 0
	case CustomOffset:
		if c.opts.CustomOffset < cp.Offset {
			cp.Offset = c.opts.CustomOffset
		}
	}
	return cp, err
}

func (c *checkpoint) Commit(offset uint64) error {
	cp := &proto.Checkpoint{
		Offset: offset,
	}
	return c.db.Put(c.key, cp)
}

func buildCheckpointKey(prefix string) []byte {
	return []byte(fmt.Sprintf("%s-checkpoint", prefix))
}
