package server

import (
	"context"

	"github.com/valli0x/auth-grpc/models"
	pb "github.com/valli0x/grpc-proto/auth-grpc/api"
)

func (s *Server) Delete(ctx context.Context, r *pb.DeleteRequest) (*pb.Empty, error) {
	resp := &pb.Empty{}

	id := r.GetId()

	//TODO: handler for admin

	entry, err := s.db.Get(ctx, id)
	if err != nil {
		return resp, err
	}

	if entry == nil {
		return resp, errNotFound
	}

	user, _ := entry.Val.(*models.User)

	// by ID
	if err := s.db.Delete(ctx, user.ID); err != nil {
		return resp, err
	}

	// by Username
	if err := s.db.Delete(ctx, user.Username); err != nil {
		return resp, err
	}

	return resp, nil
}
