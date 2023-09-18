package server

import (
	"net"

	"github.com/valli0x/auth-grpc/storage"
	pb "github.com/valli0x/grpc-proto/auth-grpc/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	UserNotFoundError = "user not found"
)

type Server struct {
	pb.UnimplementedUsersServer
	db storage.Storage
}

func NewServer(db storage.Storage) (*Server, error) {
	return &Server{
		db: db,
	}, nil
}

func (s *Server) RunServer(port string) error {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(s.ValidToken),
	}

	grpcServer := grpc.NewServer(opts...)
	reflection.Register(grpcServer)
	pb.RegisterUsersServer(grpcServer, s)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	if err := grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}
