package main

import (
	"bufio"
	"flag"
	"log"
	"os"

	"github.com/alevinval/ledger/pkg/remote"
	"google.golang.org/grpc"
)

var (
	addr     = flag.String("addr", "localhost:50051", "the address to connect to")
	writerID = flag.String("writer-id", "example-writer", "the id of this writer")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := remote.NewClient(conn)

	println("Start typing lines to be sent:")
	in := bufio.NewReader(os.Stdin)
	for {
		payload, _, _ := in.ReadLine()
		idx, err := client.Write(*writerID, payload)
		if err != nil {
			log.Fatalf("failed writing: %s", err)
			continue
		}
		log.Printf("message written (offset=%d)", idx)
	}
}
