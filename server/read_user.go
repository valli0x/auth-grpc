package server

import (
	"context"

	"github.com/valli0x/auth-grpc/models"
	pb "github.com/valli0x/grpc-proto/auth-grpc/api"
)

func (s *Server) Read(ctx context.Context, r *pb.ReadRequest) (*pb.ReadResponse, error) {
	resp := &pb.ReadResponse{
		User: &pb.User{},
	}

	id := r.GetId()

	entry, err := s.db.Get(ctx, id)
	if err != nil {
		return resp, err
	}

	if entry == nil {
		return resp, errNotFound
	}

	user, _ := entry.Val.(*models.User) // need to handle it

	resp.User.Id = user.ID
	resp.User.Username = user.Username
	resp.User.Email = user.Email

	return resp, nil
}
