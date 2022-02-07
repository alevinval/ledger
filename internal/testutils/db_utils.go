package testutils

import (
	"log"

	"github.com/dgraph-io/badger/v3"
)

// WithDB provides a badger DB instance to run tests against
func WithDB(fn func(db *badger.DB)) {
	db := mustOpenBadger()
	defer mustCloseBadger(db)

	fn(db)
}

func mustOpenBadger() *badger.DB {
	opts := badger.DefaultOptions("")
	opts.InMemory = true
	opts.Logger = nil

	db, err := badger.Open(opts)
	if err != nil {
		log.Fatalf("cannot open badger db: %s", err)
	}

	log.Printf("opened badger db (in-memory)")
	return db
}

func mustCloseBadger(db *badger.DB) {
	if !db.IsClosed() {
		err := db.Close()
		if err != nil {
			log.Fatalf("failed closing badger db: %s\n", err)
		}
	}
}
