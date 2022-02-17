package base

import "sync"

type PartitionSelector struct {
	sync.Mutex

	Size    int
	current int
}

func (ps *PartitionSelector) Next() int {
	ps.Lock()
	defer ps.Unlock()
	defer ps.increment()

	return ps.current
}

func (ps *PartitionSelector) increment() {
	ps.current++
	if ps.current == ps.Size {
		ps.current = 0
	}
}
