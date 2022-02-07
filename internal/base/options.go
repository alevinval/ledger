package base

import "time"

const (
	// LatestOffset reads from latest write offset onwards.
	LatestOffset OffsetMode = iota
	// EarliestOffset reads from earliest possible offset onwards.
	EarliestOffset
	// CustomOffset reads from a specific offset onwards.
	CustomOffset
)

type (
	// Options to configure both ledger writer and reader.
	Options struct {
		// BatchSize is used for two things:
		// - Determine the size of the intermediate buffer from which messages are read
		// - Determine the key space that is generated for efficiently scanning the kv store.
		// For the latter, the number of digits will be used to determine a power of ten (10, 100, 1000, ...)
		BatchSize    uint64
		CustomOffset uint64
		// DeliveryTimeout defines how long fetch will wait trying to queue a message to be processed
		// when the timeout is reached, the fetching is cancelled. Unit is milliseconds.
		DeliveryTimeout   time.Duration
		Offset            OffsetMode
		SequenceBandwidth uint64
	}

	// OffsetMode to determine how the initial checkpoint will be created
	OffsetMode byte
)

func (m OffsetMode) String() string {
	switch m {
	case LatestOffset:
		return "latest"
	case EarliestOffset:
		return "earliest"
	case CustomOffset:
		return "custom"
	default:
		return "unknown"
	}
}

// DefaultOptions returns most common configuration
func DefaultOptions() *Options {
	return &Options{
		BatchSize:         1000,
		DeliveryTimeout:   60000,
		Offset:            LatestOffset,
		SequenceBandwidth: 1000,
	}
}
