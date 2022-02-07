package ledger

import "github.com/alevinval/ledger/internal/base"

const (
	EarliestOffset = base.EarliestOffset
	LatestOffset   = base.LatestOffset
	CustomOffset   = base.CustomOffset
)

type OffsetMode = base.OffsetMode
type Options = base.Options

// DefaultOptions returns most common configuration
func DefaultOptions() *Options {
	return base.DefaultOptions()
}
