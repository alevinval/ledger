package ledger

import (
	"sort"

	"github.com/alevinval/ledger/internal/base"
)

type emptyObj = struct{}

var (
	empty = emptyObj{}
)

func fireAndForget(out chan emptyObj) {
	select {
	case out <- empty:
	default:
	}
}

func fireAndWait(out chan emptyObj, wait chan emptyObj) {
	out <- empty
	<-wait
}

func getSortedReaders(readers []*Reader) (sortedReaders []*Reader, anyErr error) {
	sortedReaders = make([]*Reader, len(readers))
	copy(sortedReaders, readers)

	sort.Slice(sortedReaders, func(i int, j int) bool {
		a, err := readers[i].checkpoint.GetCheckpoint()
		if err != nil {
			anyErr = err
			return false
		}
		b, err := readers[j].checkpoint.GetCheckpoint()
		if err != nil {
			anyErr = err
			return false
		}
		return a.GetOffset() < b.GetOffset()
	})

	return
}

func getChannelsForReaders(readers []*Reader) ([]<-chan base.Message, error) {
	channels := make([]<-chan base.Message, len(readers))
	for i, reader := range readers {
		ch, err := reader.Read()
		if err != nil {
			return nil, err
		}
		channels[i] = ch
	}
	return channels, nil
}
