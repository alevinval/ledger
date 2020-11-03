package testutils

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/dgraph-io/badger/v2"
)

var dbCounter int

// WithDB provides a badger DB instance to run tests against
func WithDB(fn func(db *badger.DB)) {
	db, err := openBadgerDB()
	if err != nil {
		log.Fatalf("cannot open badger db: %s", err)
	}
	defer db.Close()

	fn(db)
}

func openBadgerDB() (*badger.DB, error) {
	tmpPath := os.TempDir()
	storePath := path.Join(tmpPath, fmt.Sprintf("test-badger-%d.db", dbCounter))
	dbCounter += 1

	os.RemoveAll(storePath)

	opts := badger.DefaultOptions(storePath)
	opts.Logger = nil

	log.Printf("opening badger db in %s", storePath)
	return badger.Open(opts)
}
