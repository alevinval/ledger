package remote

import (
	"context"

	"github.com/alevinval/ledger/pkg/proto"
	"github.com/dgraph-io/badger/v3"
)

var _ proto.LedgerServer = (*Server)(nil)

type Server struct {
	proto.UnimplementedLedgerServer

	registry *Registry
	server   proto.LedgerServer
	id       string
}

func NewServer(db *badger.DB) *Server {
	return &Server{
		registry: NewRegistry(db),
	}
}

func (s *Server) Write(ctx context.Context, req *proto.WriteRequest) (*proto.WriteResponse, error) {
	writer, err := s.registry.Writer(req.WriterId)
	if err != nil {
		return nil, err
	}
	offset, err := writer.Write(req.GetData())
	return &proto.WriteResponse{
		Offset: offset,
	}, err
}

func (s *Server) Read(req *proto.ReadRequest, srv proto.Ledger_ReadServer) error {
	reader, err := s.registry.Reader(req.WriterId, req.ReaderId)
	if err != nil {
		return err
	}

	ch, err := reader.Read()
	if err != nil {
		return err
	}

	ctx := srv.Context()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case message := <-ch:
			protoMessage := &proto.Message{
				Data:   message.Data(),
				Offset: message.Offset(),
			}
			srv.Send(protoMessage)
		}
	}
}

func (s *Server) Commit(ctx context.Context, req *proto.CommitRequest) (*proto.CommitResponse, error) {
	reader, err := s.registry.Reader(req.WriterId, req.ReaderId)
	if err != nil {
		return nil, err
	}

	err = reader.Commit(req.GetOffset())
	if err != nil {
		return nil, err
	}

	return &proto.CommitResponse{}, nil
}

func (s *Server) CloseWriter(ctx context.Context, req *proto.WriterCloseRequest) (*proto.WriterCloseResponse, error) {
	s.registry.CloseWriter(req.WriterId)
	return &proto.WriterCloseResponse{}, nil
}

func (s *Server) CloseReader(ctx context.Context, req *proto.ReaderCloseRequest) (*proto.ReaderCloseResponse, error) {
	s.registry.CloseReader(req.WriterId, req.ReaderId)
	return &proto.ReaderCloseResponse{}, nil
}
