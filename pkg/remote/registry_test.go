package remote

import (
	"testing"

	"github.com/alevinval/ledger/internal/testutils"
	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/assert"
)

func TestRegistryWriter_reuseInstance(t *testing.T) {
	withRegistry(func(registry *Registry) {
		firstWriter, err := registry.Writer("writer")
		assert.NoError(t, err)

		secondWriter, err := registry.Writer("writer")
		assert.NoError(t, err)

		assert.Equal(t, firstWriter, secondWriter)
	})
}

func TestRegistryWriter_reuseDifferentInstances(t *testing.T) {
	withRegistry(func(registry *Registry) {
		firstWriter, err := registry.Writer("writer-1")
		assert.NoError(t, err)

		secondWriter, err := registry.Writer("writer-2")
		assert.NoError(t, err)

		assert.NotEqual(t, firstWriter, secondWriter)
	})
}

func TestRegistryReader_whenMissingWriter_returnsErr(t *testing.T) {
	withRegistry(func(registry *Registry) {
		_, err := registry.Reader("missing-writer", "reader-1")
		assert.ErrorIs(t, ErrWriterNotFound, err)
	})
}

func TestRegistryReader_whenWriterAvailable_reusesInstance(t *testing.T) {
	withRegistry(func(registry *Registry) {
		_, err := registry.Writer("writer")
		assert.NoError(t, err)

		firstReader, err := registry.Reader("writer", "reader-1")
		assert.NoError(t, err)

		secondReader, err := registry.Reader("writer", "reader-1")
		assert.NoError(t, err)

		assert.Equal(t, firstReader, secondReader)
	})
}

func TestRegistryReader_whenWriterAvailable_differentInstances(t *testing.T) {
	withRegistry(func(registry *Registry) {
		_, err := registry.Writer("writer")
		assert.NoError(t, err)

		firstReader, err := registry.Reader("writer", "reader-1")
		assert.NoError(t, err)

		secondReader, err := registry.Reader("writer", "reader-2")
		assert.NoError(t, err)

		assert.NotEqual(t, firstReader, secondReader)
	})
}

func withRegistry(fn func(registry *Registry)) {
	testutils.WithDB(func(db *badger.DB) {
		registry := NewRegistry(db)
		defer registry.Close()

		fn(registry)
	})
}
