package storage

import (
	"testing"

	"github.com/alevinval/ledger/internal/base"
	"github.com/alevinval/ledger/internal/testutils"
	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/assert"
)

var opts = base.DefaultOptions()

func TestStorageGetBytes(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		s := &Storage{db, opts}

		_, err := s.GetBytes([]byte("some-key"))

		assert.Error(t, err)
	})
}

func TestStoragePutBytes(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		s := &Storage{db, opts}

		err := s.PutBytes([]byte("some-key"), []byte("some-value"))

		assert.NoError(t, err)
	})
}

func TestStoragePutAndGetBytes(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		s := &Storage{db, opts}

		s.PutBytes([]byte("some-key"), []byte("some-value"))
		value, err := s.GetBytes([]byte("some-key"))

		assert.NoError(t, err)
		assert.Equal(t, []byte("some-value"), value)
	})
}

func TestStorageScanKeysIndexedEmpty(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		s := &Storage{db, opts}

		keys := [][]byte{}
		s.ScanKeysIndexed([]byte("some-prefix-"), 0, func(k []byte, idx uint64) error {
			keys = append(keys, k)
			return nil
		})

		assert.Equal(t, 0, len(keys))
	})
}

func TestStorageScanKeysIndexed(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		s := &Storage{db, opts}

		s.PutBytes([]byte("prefix-000000000000000001"), []byte("value-1"))
		s.PutBytes([]byte("prefix-000000000000000002"), []byte("value-2"))
		s.PutBytes([]byte("prefix-000000000000000003"), []byte("value-3"))
		s.PutBytes([]byte("prefix-000000000000000004"), []byte("value-4"))
		s.PutBytes([]byte("prefix-000000000000000005"), []byte("value-5"))

		keys := [][]byte{}
		s.ScanKeysIndexed([]byte("prefix-"), 2, func(k []byte, idx uint64) error {
			keys = append(keys, k)
			return nil
		})

		assert.Equal(
			t,
			[][]byte{
				[]byte("prefix-000000000000000003"),
				[]byte("prefix-000000000000000004"),
				[]byte("prefix-000000000000000005"),
			},
			keys)
	})
}
