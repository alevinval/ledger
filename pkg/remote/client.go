package remote

import (
	"context"
	"io"
	"log"

	"github.com/alevinval/ledger/pkg/proto"
	"google.golang.org/grpc"
)

type Client struct {
	proto.LedgerClient
}

func NewClient(conn *grpc.ClientConn) *Client {
	return &Client{
		proto.NewLedgerClient(conn),
	}
}

func (c *Client) Write(writerID string, data []byte) (uint64, error) {
	writeReq := &proto.WriteRequest{
		WriterId: writerID,
		Data:     data,
	}
	writeRes, err := c.LedgerClient.Write(context.Background(), writeReq)
	if err != nil {
		return 0, err
	}
	return writeRes.Offset, nil
}

func (c *Client) Read(writerID, readerID string) (<-chan *proto.Message, error) {
	readReq := &proto.ReadRequest{
		WriterId: writerID,
		ReaderId: readerID,
	}
	stream, err := c.LedgerClient.Read(context.Background(), readReq)
	if err != nil {
		return nil, err
	}

	out := make(chan *proto.Message)

	go func() {
		for {
			message, err := stream.Recv()
			if err == io.EOF {
				close(out)
				return
			}

			if err != nil {
				log.Fatalf("cannot receive %v", err)
			}
			out <- message
		}
	}()

	return out, nil
}

func (c *Client) Commit(writerID, readerID string, offset uint64) (*proto.CommitResponse, error) {
	commitReq := &proto.CommitRequest{
		WriterId: writerID,
		ReaderId: readerID,
		Offset:   offset,
	}
	return c.LedgerClient.Commit(context.Background(), commitReq)
}

func (c *Client) CloseWriter(writerID string) {
	closeReq := &proto.WriterCloseRequest{
		WriterId: writerID,
	}
	c.LedgerClient.CloseWriter(context.Background(), closeReq)
}

func (c *Client) CloseReader(writerID, readerID string) {
	closeReq := &proto.ReaderCloseRequest{
		WriterId: writerID,
		ReaderId: readerID,
	}
	c.LedgerClient.CloseReader(context.Background(), closeReq)
}
