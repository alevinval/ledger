package remote

import (
	"fmt"
	"net"
	"testing"

	"github.com/alevinval/ledger/internal/testutils"
	"github.com/alevinval/ledger/pkg/proto"
	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestServerAndClient(t *testing.T) {
	WithLedgerServer(t, func(client *Client) {
		rsp, err := client.Write("some-channel", []byte("hello world"))
		assert.NoError(t, err)
		assert.Equal(t, uint64(1), rsp)
	})
}

func WithLedgerServer(t *testing.T, test func(client *Client)) {
	testutils.WithDB(func(db *badger.DB) {
		listener, err := net.Listen("tcp", getAddress(50051))
		if err != nil {
			t.Fatalf("ledger server cannot listen: %s", err)
		}

		grpcServer := grpc.NewServer()
		ledgerServer := NewServer(db)
		proto.RegisterLedgerServer(grpcServer, ledgerServer)

		go grpcServer.Serve(listener)
		defer grpcServer.Stop()

		clientConn, err := grpc.Dial(getAddress(50051), grpc.WithInsecure())
		assert.NoError(t, err)
		defer clientConn.Close()

		ledgerClient := NewClient(clientConn)
		test(ledgerClient)
	})
}

func getAddress(port int) string {
	return fmt.Sprintf("localhost:%d", port)
}
