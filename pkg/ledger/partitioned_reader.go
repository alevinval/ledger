package ledger

import (
	"github.com/alevinval/ledger/internal/base"
	"go.uber.org/zap"
)

// PartitionedReader reads using a partition scheme
type PartitionedReader struct {
	id         string
	partitions int
	out        chan base.PartitionedMessage
	readers    []*Reader
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
		out:        make(chan base.PartitionedMessage, 1),
	}

	err = r.startFetcher()
	if err != nil {
		logger.Error("cannot start fetcher", zap.String("id", readerID), zap.Error(err))
		return nil, err
	}

	return r, nil
}

func (r *PartitionedReader) Read() (<-chan base.PartitionedMessage, error) {
	return r.out, nil
}

func (r *PartitionedReader) Close() {
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

func (r *PartitionedReader) fetcher(channels []<-chan base.Message) {
	for {
		for i := range channels {
			msg, open := <-channels[i]
			if !open {
				logger.Debug("stopping partitioned reader fetcher", zap.String("id", r.id))
				return
			}
			r.out <- &partitionedMessageImpl{
				partition: r.readers[i].checkpoint,
				offset:    msg.Offset(),
				data:      msg.Data(),
			}
		}
	}
}
