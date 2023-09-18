package server

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/valli0x/auth-grpc/models"
	"github.com/valli0x/auth-grpc/storage"
	pb "github.com/valli0x/grpc-proto/auth-grpc/api"
)

func (s *Server) Create(ctx context.Context, r *pb.CreateRequest) (*pb.CreateResponse, error) {
	resp := &pb.CreateResponse{}

	email := r.User.GetEmail()
	username := r.User.GetUsername()
	password := r.User.GetPassword()

	if s.exits(ctx, username) {
		return resp, errAlreadyExist
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return resp, err
	}

	user := &models.User{
		ID:       uuid,
		Email:    email,
		Username: username,
		Password: []byte(password), // maybe sha256 hash?
	}

	// by ID
	if err := s.db.Put(ctx, &storage.Entry{
		Key: user.ID,
		Val: user,
	}); err != nil {
		return resp, err
	}

	// by Username
	if err := s.db.Put(ctx, &storage.Entry{
		Key: user.Username,
		Val: user,
	}); err != nil {
		return resp, err
	}

	resp.Id = user.ID
	return resp, nil
}
