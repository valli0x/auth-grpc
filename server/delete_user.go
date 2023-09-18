package server

import (
	"context"

	pb "github.com/valli0x/grpc-proto/auth-grpc/api"
)

func (s *Server) Delete(ctx context.Context, r *pb.DeleteRequest) (*pb.Empty, error) {
	resp := &pb.Empty{}

	id := r.GetId()

	if err := s.db.Delete(ctx, id); err != nil {
		return resp, err
	}

	return resp, nil
}
