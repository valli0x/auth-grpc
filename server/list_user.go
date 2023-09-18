package server

import (
	"context"

	"github.com/valli0x/auth-grpc/models"
	pb "github.com/valli0x/grpc-proto/auth-grpc/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) List(ctx context.Context, r *pb.Empty) (*pb.ListResponse, error) {
	resp := &pb.ListResponse{
		Users: []*pb.User{},
	}

	ids, err := s.db.List(ctx)
	if err != nil {
		return resp, status.Errorf(codes.Internal, err.Error())
	}

	usernames := map[string]struct{}{}
	for _, id := range ids {
		entry, err := s.db.Get(ctx, id)
		if err != nil {
			return resp, err
		}

		user, ok := entry.Val.(*models.User)
		if !ok {
			return resp, status.Error(codes.Internal, "")
		}

		if _, ok := usernames[user.ID]; ok {
			continue
		}

		resp.Users = append(resp.Users, &pb.User{
			Id:       user.ID,
			Email:    user.Email,
			Username: user.Username,
		})
		usernames[user.ID] = struct{}{}
	}

	return resp, nil
}
