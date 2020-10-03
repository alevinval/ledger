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

func (c *checkpoint) GetCheckpoint() (*proto.CheckPoint, error) {
	var cp = &proto.CheckPoint{}
	err := c.db.Get(c.key, cp)
	return cp, err
}

func (c *checkpoint) GetCheckpointFrom(other *checkpoint) (*proto.CheckPoint, error) {
	cp, err := other.GetCheckpoint()
	switch c.opts.Mode {
	case ModeEarliest:
		cp.Index = 0
	case ModeCustom:
		if c.opts.CustomIndex < cp.Index {
			cp.Index = c.opts.CustomIndex
		}
	}
	return cp, err
}

func (c *checkpoint) Commit(index uint64) error {
	cp := &proto.CheckPoint{
		Index: index,
	}
	return c.db.Put(c.key, cp)
}

func buildCheckpointKey(prefix string) []byte {
	return []byte(fmt.Sprintf("%s-checkpoint", prefix))
}
