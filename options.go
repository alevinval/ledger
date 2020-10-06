package ledger

const (
	ModeEarliest OptionMode = iota
	ModeLatest
	ModeCustom
)

type (
	// Options to configure both ledger writer and reader.
	Options struct {
		Mode              OptionMode
		CustomIndex       uint64
		SequenceBandwidth uint64

		// BatchSize is used for two things:
		// - Determine the size of the intermediate buffer from which messages are read
		// - Determine the key space that is generated for efficiently scanning the kv store.
		// For the latter, the number of digits will be used to determine a power of ten (10, 100, 1000, ...)
		BatchSize uint64
	}

	// OptionMode to determine how the initial checkpoint will be created
	OptionMode byte
)

func DefaultOptions() *Options {
	return &Options{
		Mode:              ModeLatest,
		SequenceBandwidth: 1000,
		BatchSize:         1000,
	}
}
