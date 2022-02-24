package ledger

import (
	"github.com/alevinval/ledger/internal/base"
	"github.com/alevinval/ledger/internal/checkpoint"
)

type Message = base.BaseMessage

var _ Message = (*messageImpl)(nil)

type (
	messageImpl struct {
		chk    *checkpoint.Checkpoint
		data   []byte
		offset uint64
	}
)

func (m *messageImpl) Offset() uint64 {
	return m.offset
}

func (m *messageImpl) Data() []byte {
	return m.data
}

func (m *messageImpl) Commit() error {
	return m.chk.Commit(m.offset)
}
