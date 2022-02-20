package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/alevinval/ledger/internal/testutils"
	"github.com/alevinval/ledger/pkg/proto"
	"github.com/alevinval/ledger/pkg/remote"
	"github.com/dgraph-io/badger/v3"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	flag.Parse()
	testutils.WithDB(func(db *badger.DB) {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
		if err != nil {
			log.Fatalf("error listening: %s", err)
		}

		grpcServer := grpc.NewServer()
		log.Printf("server listening at %s", listener.Addr())

		server := remote.NewServer(db)
		proto.RegisterLedgerServer(grpcServer, server)

		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	})
}
