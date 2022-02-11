package ledger

import (
	"github.com/alevinval/ledger/internal/base"
	"github.com/alevinval/ledger/internal/checkpoint"
)

var _ base.Message = (*messageImpl)(nil)
var _ base.PartitionedMessage = (*partitionedMessageImpl)(nil)

type (
	messageImpl struct {
		offset uint64
		data   []byte
	}
	partitionedMessageImpl struct {
		partition *checkpoint.Checkpoint
		offset    uint64
		data      []byte
	}
)

func (m *messageImpl) Offset() uint64 {
	return m.offset
}

func (m *messageImpl) Data() []byte {
	return m.data
}

func (pm *partitionedMessageImpl) Commit() error {
	return pm.partition.Commit(pm.offset)
}

func (pm *partitionedMessageImpl) Offset() uint64 {
	return pm.offset
}

func (pm *partitionedMessageImpl) Data() []byte {
	return pm.data
}
