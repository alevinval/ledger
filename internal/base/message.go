package base

type (
	// Message interface used to pass read results from a ledger.Reader
	// Provides the data payload and the offset in the log
	Message interface {
		Offset() uint64
		Data() []byte
	}

	// PartitionedMessage interface used to pass read results from a ledger.PartitionedReader
	// Provides the data payload and the offset in the log.
	//
	// Additionally, provides an interface to commit the reads, when using a PartitionedReader
	// the partition from which the message is read is hidden from the user, this simplifies
	// the usage of the library.
	PartitionedMessage interface {
		Commit() error
		Offset() uint64
		Data() []byte
	}
)
