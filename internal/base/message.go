package base

type (
	// Message interface used to pass read results from a ledger.Reader
	// Provides the data payload and the offset in the log
	Message interface {
		Commit() error
		Data() []byte
		Offset() uint64
	}
)
