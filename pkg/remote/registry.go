package remote

import (
	"errors"
	"fmt"
	"sync"

	"github.com/alevinval/ledger/internal/base"
	"github.com/alevinval/ledger/pkg/ledger"
	"github.com/dgraph-io/badger/v3"
)

var ErrWriterNotFound = errors.New("writer not found")

type Registry struct {
	mu sync.RWMutex

	db      *badger.DB
	writers map[string]*ledger.Writer
	readers map[string]*ledger.Reader
}

func NewRegistry(db *badger.DB) *Registry {
	return &Registry{
		db:      db,
		writers: make(map[string]*ledger.Writer),
		readers: make(map[string]*ledger.Reader),
	}
}

func (r *Registry) Writer(id string) (*ledger.Writer, error) {
	return r.WriterOpts(id, base.DefaultOptions())
}

func (r *Registry) WriterOpts(id string, opts *base.Options) (writer *ledger.Writer, err error) {
	writer, ok := r.writers[id]
	if ok {
		return writer, nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	writerKey := getWriterKey(id)
	writer, ok = r.writers[writerKey]
	if ok {
		return writer, nil
	}

	writer, err = ledger.NewWriterOpts(id, r.db, opts)
	if err != nil {
		return nil, err
	}
	r.writers[writerKey] = writer
	return writer, nil
}

func (r *Registry) Reader(writerId, readerId string) (*ledger.Reader, error) {
	opts := base.DefaultOptions()
	opts.Offset = base.CustomOffset
	opts.CustomOffset = 0
	return r.ReaderOpts(writerId, readerId, opts)
}

func (r *Registry) ReaderOpts(writerId, readerId string, opts *base.Options) (reader *ledger.Reader, err error) {
	readerKey := getReaderKey(writerId, readerId)
	reader, ok := r.readers[readerKey]
	if ok {
		return reader, nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	reader, ok = r.readers[readerKey]
	if ok {
		return reader, nil
	}

	writerKey := getWriterKey(writerId)
	writer, ok := r.writers[writerKey]
	if !ok {
		return nil, ErrWriterNotFound
	}

	reader, err = writer.NewReaderOpts(readerId, opts)
	if err != nil {
		return nil, err
	}
	r.readers[readerKey] = reader
	return reader, err
}

func (r *Registry) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, writer := range r.writers {
		writer.Close()
	}

	for _, reader := range r.readers {
		reader.Close()
	}
}

func (r *Registry) CloseWriter(writerId string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	writerKey := getWriterKey(writerId)
	writer, ok := r.writers[writerKey]
	if ok {
		writer.Close()
		delete(r.writers, writerKey)
	}
}

func (r *Registry) CloseReader(writerId string, readerId string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	readerKey := getReaderKey(writerId, readerId)
	reader, ok := r.readers[readerKey]
	if ok {
		reader.Close()
		delete(r.readers, readerKey)
	}
}

func getWriterKey(writerId string) string {
	return fmt.Sprintf("w[%s]", writerId)
}

func getReaderKey(writerId, readerId string) string {
	return fmt.Sprintf("w[%s]r[%s]", writerId, readerId)
}
