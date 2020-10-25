package ledger

import (
	"bytes"
	"fmt"
	"math"
	"strconv"

	"github.com/dgraph-io/badger/v2"
	"github.com/golang/protobuf/proto"
)

const maxOffset = math.MaxUint64

var (
	uintDigits = getDigitsFromNumber(maxOffset)
)

// storage wrapper for all badger related operations, mostly for convenience
// but also critical pieces the ledger relies upon, like the ScanKeysIndexed
// implementation which makes everything possible.
type storage struct {
	db   *badger.DB
	opts *Options
}

func (s *storage) GetBytes(key []byte) (value []byte, err error) {
	err = s.db.View(func(txn *badger.Txn) error {
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

func (s *storage) PutBytes(key, value []byte) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

func (s *storage) Get(key []byte, dst proto.Message) error {
	return s.db.View(func(txn *badger.Txn) error {
		entry, err := txn.Get(key)
		if err != nil {
			return err
		}
		return entry.Value(func(value []byte) error {
			return proto.Unmarshal(value, dst)
		})
	})
}

func (s *storage) Put(key []byte, pb proto.Message) error {
	value, err := proto.Marshal(pb)
	if err != nil {
		return err
	}
	return s.PutBytes(key, value)
}

func (s *storage) ScanKeysIndexed(
	basePrefix []byte,
	startOffset uint64,
	fn func(key []byte, offset uint64) (err error),
) error {
	return s.db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		for prefix := range s.buildKeySpace(basePrefix, startOffset) {
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

func (s *storage) buildKeySpace(prefix []byte, startOffset uint64) <-chan []byte {
	out := make(chan []byte)
	go func() {
		keySpaceFmt := getKeySpaceScanFmt(s.opts)
		for offset := startOffset / s.opts.BatchSize; offset < maxOffset; offset++ {
			prefix := fmt.Sprintf(keySpaceFmt, prefix, offset)
			out <- []byte(prefix)
			logger.Log("storage", "buildKeySpace", "generatedKey", prefix)
		}
	}()
	return out
}

func getKeySpaceScanFmt(opts *Options) string {
	batchDigits := getDigitsFromNumber(opts.BatchSize)
	return "%s%0." + strconv.Itoa(uintDigits-batchDigits) + "d"
}

func getDigitsFromNumber(number uint64) (numDigits int) {
	for i := number; i >= 10; i /= 10 {
		numDigits++
	}
	if numDigits == 0 {
		return 1
	}
	return
}
