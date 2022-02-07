package storage

import (
	"testing"

	"github.com/alevinval/ledger/internal/base"
	"github.com/alevinval/ledger/internal/testutils"
	"github.com/alevinval/ledger/pkg/proto"
	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/assert"
)

var defaultOpts = base.DefaultOptions()

func TestStorageGetBytes_missingKey(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		s := &Storage{db, defaultOpts}

		_, err := s.GetBytes([]byte("missing-key"))

		assert.ErrorIs(t, err, badger.ErrKeyNotFound)
	})
}

func TestStoragePutBytes(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		s := &Storage{db, defaultOpts}

		err := s.PutBytes([]byte("some-key"), []byte("some-value"))

		assert.NoError(t, err)
	})
}

func TestStoragePutBytes_thenGetBytes(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		s := &Storage{db, defaultOpts}

		err := s.PutBytes([]byte("some-key"), []byte("some-value"))
		assert.NoError(t, err)

		value, err := s.GetBytes([]byte("some-key"))
		assert.NoError(t, err)
		assert.Equal(t, []byte("some-value"), value)
	})
}

func TestStorageScanKeysIndexed_emptyResult(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		s := &Storage{db, defaultOpts}

		keys := [][]byte{}
		s.ScanKeysIndexed([]byte("some-prefix-"), 0, func(k []byte, idx uint64) error {
			keys = append(keys, k)
			return nil
		})

		assert.Equal(t, 0, len(keys))
	})
}

func TestStorageScanKeysIndexed_matchingResultsFromIndex(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		s := &Storage{db, defaultOpts}

		startFrom := 2
		s.PutBytes([]byte("prefix-000000000000000001"), []byte("value-1"))
		s.PutBytes([]byte("prefix-000000000000000002"), []byte("value-2"))
		s.PutBytes([]byte("prefix-000000000000000003"), []byte("value-3"))
		s.PutBytes([]byte("prefix-000000000000000004"), []byte("value-4"))
		s.PutBytes([]byte("prefix-000000000000000005"), []byte("value-5"))

		keys := [][]byte{}
		s.ScanKeysIndexed([]byte("prefix-"), uint64(startFrom), func(k []byte, idx uint64) error {
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

func TestStorageGet_missingKey(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		s := &Storage{db, defaultOpts}
		dst := &proto.Checkpoint{}

		err := s.Get([]byte("missing-key"), dst)

		assert.ErrorIs(t, err, badger.ErrKeyNotFound)
	})
}

func TestStoragePut(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		s := &Storage{db, defaultOpts}
		obj := &proto.Checkpoint{
			Offset: 123,
		}

		err := s.Put([]byte("some-key"), obj)

		assert.NoError(t, err)
	})
}

func TestStoragePut_thenGet(t *testing.T) {
	testutils.WithDB(func(db *badger.DB) {
		s := &Storage{db, defaultOpts}
		expected := &proto.Checkpoint{
			Offset: 123,
		}

		err := s.Put([]byte("some-key"), expected)
		assert.NoError(t, err)

		actual := &proto.Checkpoint{}
		err = s.Get([]byte("some-key"), actual)
		assert.NoError(t, err)

		assert.Equal(t, expected.Offset, actual.Offset)
	})
}
