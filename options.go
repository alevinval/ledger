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

		// KeySpaceBatchSize must be apower of ten (10, 100, 1000, 10000 ...)
		// When set to something else, the number of digits will determine which
		// power is used.
		KeySpaceBatchSize uint64
	}

	// OptionMode to determine how the initial checkpoint will be created
	OptionMode byte
)

func DefaultOptions() *Options {
	return &Options{
		Mode:              ModeLatest,
		SequenceBandwidth: 1000,
		KeySpaceBatchSize: 1000,
	}
}
