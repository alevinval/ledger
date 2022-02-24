package base

type (
	// BaseMessage interface used to pass read results from a ledger.Reader
	// Provides the data payload and the offset in the log
	BaseMessage interface {
		Commit() error
		Data() []byte
		Offset() uint64
	}
)
