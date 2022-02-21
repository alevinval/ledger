package exp

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
	r.mu.Lock()
	defer r.mu.Unlock()

	writer, ok := r.writers[id]
	if ok {
		return writer, nil
	}

	writer, err = ledger.NewWriterOpts(id, r.db, opts)
	if err != nil {
		return nil, err
	}
	r.writers[id] = writer
	return writer, nil
}

func (r *Registry) Reader(writerId, readerId string) (*ledger.Reader, error) {
	return r.ReaderOpts(writerId, readerId, base.DefaultOptions())
}

func (r *Registry) ReaderOpts(writerId, readerId string, opts *base.Options) (reader *ledger.Reader, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := fmt.Sprintf("r[%s]w[%s]", readerId, writerId)
	reader, ok := r.readers[key]
	if ok {
		return reader, nil
	}

	writer, ok := r.writers[writerId]
	if !ok {
		return nil, ErrWriterNotFound
	}

	reader, err = writer.NewReaderOpts(readerId, opts)
	if err != nil {
		return nil, err
	}
	r.readers[key] = reader
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
