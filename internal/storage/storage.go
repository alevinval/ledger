package storage

import (
	"bytes"
	"fmt"
	"math"
	"strconv"

	"github.com/alevinval/ledger/internal/base"
	"github.com/alevinval/ledger/internal/util/log"
	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const maxOffset = math.MaxUint64

var (
	UintDigits = getDigitsFromNumber(maxOffset)

	logger = log.GetLogger()
)

// storage wrapper for all badger related operations, mostly for convenience
// but also critical pieces the ledger relies upon, like the ScanKeysIndexed
// implementation which makes everything possible.
type Storage struct {
	DB   *badger.DB
	Opts *base.Options
}

type keySpaceIterator struct {
	opts       *base.Options
	basePrefix []byte
	offset     uint64
}

func (s *Storage) GetBytes(key []byte) (value []byte, err error) {
	err = s.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		return item.Value(func(v []byte) error {
			value = v
			return nil
		})
	})
	return
}

func (s *Storage) PutBytes(key, value []byte) error {
	return s.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

func (s *Storage) Get(key []byte, dst proto.Message) error {
	return s.DB.View(func(txn *badger.Txn) error {
		entry, err := txn.Get(key)
		if err != nil {
			return err
		}
		return entry.Value(func(value []byte) error {
			return proto.Unmarshal(value, dst)
		})
	})
}

func (s *Storage) Put(key []byte, pb proto.Message) error {
	value, err := proto.Marshal(pb)
	if err != nil {
		return err
	}
	return s.PutBytes(key, value)
}

func (s *Storage) ScanKeysIndexed(
	basePrefix []byte,
	startOffset uint64,
	fn func(key []byte, offset uint64) (err error),
) error {
	return s.DB.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		keySpaceIt := s.newKeySpaceIterator(basePrefix, startOffset)
		for keySpaceIt.HasNext() {
			prefix := keySpaceIt.Next()
			isEmptySeek := true
			it.Rewind()
			for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
				item := it.Item()
				keyOffset := bytes.TrimPrefix(item.Key(), basePrefix)
				offset, err := strconv.ParseUint(string(keyOffset), 10, 64)
				if err != nil {
					continue
				}
				if offset > startOffset {
					err := fn(item.KeyCopy(nil), offset)
					if err != nil {
						return err
					}
					isEmptySeek = false
				}
			}

			// When an empty seek occurs abort the scanning since the end
			// of the indexed records has been reached. This assumes that
			// there are no gaps between writes.
			if isEmptySeek {
				break
			}
		}
		return nil
	})
}

func (s *Storage) newKeySpaceIterator(prefix []byte, startOffset uint64) *keySpaceIterator {
	return &keySpaceIterator{
		opts:       s.Opts,
		basePrefix: prefix,
		offset:     startOffset / s.Opts.BatchSize,
	}
}

func (it *keySpaceIterator) HasNext() bool {
	return it.offset < maxOffset
}

func (it *keySpaceIterator) Next() []byte {
	keySpaceFmt := getKeySpaceScanFmt(it.opts.BatchSize)
	prefix := fmt.Sprintf(keySpaceFmt, it.basePrefix, it.offset)
	it.offset++
	logger.Debug("storage.keySpaceIterator.Next", zap.String("prefix", prefix))
	return []byte(prefix)
}

func getKeySpaceScanFmt(batchSize uint64) string {
	batchDigits := getDigitsFromNumber(batchSize)
	return "%s%0." + strconv.Itoa(UintDigits-batchDigits) + "d"
}

func getDigitsFromNumber(number uint64) (numDigits int) {
	for number > 1 {
		if number > 1e16 {
			numDigits += 16
			number /= 1e16
		}
		if number >= 1e8 {
			numDigits += 8
			number /= 1e8
		}
		if number >= 1e4 {
			numDigits += 4
			number /= 1e4
		}
		if number >= 1e2 {
			numDigits += 2
			number /= 1e2
		}
		if number >= 1e1 {
			numDigits += 1
			number /= 1e1
		}
	}
	numDigits++
	return
}
