package ledger

import "fmt"

// PartitionedReader reads using a partition scheme
type PartitionedReader struct {
	readers    []*Reader
	out        chan *PartitionedMessage
	partitions int
	current    int
}

func (w *PartitionedWriter) NewReader(id string) (*PartitionedReader, error) {
	readers := make([]*Reader, w.partitions)
	for i := range w.writers {
		r, err := w.writers[i].NewReader(fmt.Sprintf("%s-part-%d", id, i))
		if err != nil {
			return nil, err
		}
		readers[i] = r
	}

	r := &PartitionedReader{
		readers:    readers,
		partitions: w.partitions,
		out:        make(chan *PartitionedMessage, 1),
	}

	go r.fetcher()

	return r, nil
}

func (w *PartitionedReader) Read() (<-chan *PartitionedMessage, error) {
	return w.out, nil
}

func (r *PartitionedReader) fetcher() error {
	channels := make([]<-chan *Message, r.partitions)
	for i := range r.readers {
		ch, err := r.readers[i].Read()
		if err != nil {
			return err
		}
		channels[i] = ch
	}

	for {
		for i := range channels {
			msg := <-channels[i]
			r.out <- &PartitionedMessage{
				Partition: r.readers[i].chk,
				Offset:    msg.Offset,
				Data:      msg.Data,
			}
		}
	}
}
