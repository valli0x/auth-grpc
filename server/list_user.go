package server

import (
	"context"

	"github.com/valli0x/auth-grpc/models"
	pb "github.com/valli0x/grpc-proto/auth-grpc/api"
)

func (s *Server) List(ctx context.Context, r *pb.Empty) (*pb.ListResponse, error) {
	resp := &pb.ListResponse{
		Users: []*pb.User{},
	}

	ids, err := s.db.List(ctx)
	if err != nil {
		return resp, err
	}

	for _, id := range ids {
		entry, err := s.db.Get(ctx, id)
		if err != nil {
			return resp, err
		}

		user, ok := entry.Val.(*models.User)
		if !ok {
			return resp, nil
		}

		resp.Users = append(resp.Users, &pb.User{
			Id:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		})
	}

	return resp, nil
}
