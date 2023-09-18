package server

import (
	"context"

	pb "github.com/valli0x/grpc-proto/auth-grpc/api"
)

func (s *Server) Delete(ctx context.Context, r *pb.DeleteRequest) (*pb.Empty, error) {
	id := r.User.GetId()

	if err := s.db.Delete(ctx, id); err != nil {
		return &pb.Empty{}, err
	}

	return &pb.Empty{}, nil
}
