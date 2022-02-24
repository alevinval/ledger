package ledger

import (
	"sync"

	"go.uber.org/zap"
)

// PartitionedReader reads using a partition scheme
type PartitionedReader struct {
	mu sync.RWMutex

	out        chan Message
	id         string
	readers    []*Reader
	partitions int

	isClosed bool
}

func (pw *PartitionedWriter) NewReader(readerID string) (*PartitionedReader, error) {
	readers, err := pw.createReaders(readerID)
	if err != nil {
		logger.Error("failed creating readers for partitioned writer", zap.String("id", readerID), zap.Error(err))
		return nil, err
	}

	r := &PartitionedReader{
		id:         readerID,
		readers:    readers,
		partitions: pw.partitions,
		out:        make(chan Message, 1),
	}

	err = r.startFetcher()
	if err != nil {
		logger.Error("cannot start fetcher", zap.String("id", readerID), zap.Error(err))
		return nil, err
	}

	return r, nil
}

func (r *PartitionedReader) Read() (<-chan Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.isClosed {
		return nil, ErrClosedReader
	}

	return r.out, nil
}

func (r *PartitionedReader) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.isClosed = true

	for i := range r.readers {
		r.readers[i].Close()
	}
}

func (r *PartitionedReader) startFetcher() error {
	sortedReaders, err := getSortedReaders(r.readers)
	if err != nil {
		logger.Error("fetcher error: cannot get sorted readers", zap.String("id", r.id), zap.Error(err))
		return err
	}

	channels, err := getChannelsForReaders(sortedReaders)
	if err != nil {
		logger.Error("fetcher error: cannot get reader channels", zap.String("id", r.id), zap.Error(err))
		return err
	}

	go r.fetcher(channels)

	return nil
}

func (r *PartitionedReader) fetcher(channels []<-chan Message) {
	for {
		for i := range channels {
			msg, open := <-channels[i]
			if !open {
				logger.Debug("stopping partitioned reader fetcher", zap.String("id", r.id))
				return
			}
			r.out <- &messageImpl{
				chk:    r.readers[i].checkpoint,
				data:   msg.Data(),
				offset: msg.Offset(),
			}
		}
	}
}
