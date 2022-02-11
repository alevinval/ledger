package base

type Commitable interface {
	Commit(uint64) error
}

type (
	// Message structure used to represent read results
	Message interface {
		Offset() uint64
		Data() []byte
	}

	// PartitionedMessage structure that represents messages from a partitioned reader
	PartitionedMessage interface {
		Commit() error
		Offset() uint64
		Data() []byte
	}
)
