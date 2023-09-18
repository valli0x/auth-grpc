package server

import (
	"github.com/valli0x/auth-grpc/storage"
	pb "github.com/valli0x/grpc-proto/auth-grpc/api"
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
