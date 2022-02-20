package main

import (
	"flag"
	"log"

	"github.com/alevinval/ledger/pkg/remote"
	"google.golang.org/grpc"
)

var (
	addr     = flag.String("addr", "localhost:50051", "the address to connect to")
	writerID = flag.String("writer-id", "example-writer", "id of the writer to read from")
	readerID = flag.String("reader-id", "example-reader", "id of this reader")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := remote.NewClient(conn)

	messages, err := client.Read(*writerID, *readerID)
	if err != nil {
		log.Fatalf("could not read: %v", err)
	}

	for msg := range messages {
		log.Printf("read message: %s (offset=%d)", msg.GetData(), msg.GetOffset())
		_, err := client.Commit(*writerID, *readerID, msg.GetOffset())
		if err != nil {
			log.Printf("error committing offset: %s", err)
			break
		}
	}
}
