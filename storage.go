package ledger

import (
	"bytes"
	"fmt"
	"math"
	"strconv"

	"github.com/dgraph-io/badger"
	"github.com/golang/protobuf/proto"
)

const maxIdx = math.MaxUint64

var (
	uintDigits int
)

func init() {
	var i uint64
	for i = maxIdx; i >= 10; i /= 10 {
		uintDigits++
	}
}

// storage wrapper for all badger related operations, mostly for convenience
// but also critical pieces the ledger relies upon, like the ScanKeysIndexed
// implementation which makes everything possible.
type storage struct {
	db   *badger.DB
	opts *Options
}

func (s *storage) PutBytes(key, value []byte) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

func (s *storage) GetBytes(key []byte) (value []byte, err error) {
	err = s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return nil
		}
		return item.Value(func(v []byte) error {
			value = v
			return nil
		})
	})
	return
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

func (s *storage) ScanKeysIndexed(basePrefix []byte, startIdx uint64, cb func(k []byte, idx uint64) (err error)) error {
	return s.db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		for prefix := range s.buildKeySpace(basePrefix, startIdx) {
			it.Rewind()

			isEmptySeek := true
			for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
				k := it.Item().KeyCopy(nil)
				idx, _ := strconv.ParseUint(string(bytes.TrimPrefix(k, prefix)), 10, 64)
				if idx > startIdx {
					err := cb(k, idx)
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

func (s *storage) buildKeySpace(prefix []byte, startIdx uint64) <-chan []byte {
	out := make(chan []byte)
	go func() {
		keySpaceFmt := getKeySpaceScanFmt(s.opts)
		for n := startIdx / s.opts.KeySpaceBatchSize; n < maxIdx; n++ {
			prefix := fmt.Sprintf(keySpaceFmt, prefix, n)
			logger.Log("generatedKey", prefix)
			out <- []byte(prefix)
		}
	}()
	return out
}

func getKeySpaceScanFmt(opts *Options) string {
	batchDigits := 0
	for i := opts.KeySpaceBatchSize; i >= 10; i /= 10 {
		batchDigits++
	}
	if batchDigits == 0 {
		batchDigits = 1
	}
	return "%s-%0." + strconv.Itoa(uintDigits-batchDigits) + "d"
}
