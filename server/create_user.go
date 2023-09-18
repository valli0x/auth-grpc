package server

import (
	"context"
	"errors"

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

	if s.exits(username) {
		return resp, errors.New("username is not unique")
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:       uuid,
		Email:    email,
		Username: username,
		Password: []byte(password), // maybe sha256 hash?
	}

	if err := s.db.Put(ctx, &storage.Entry{
		Key: user.ID,
		Val: user,
	}); err != nil {
		return nil, err
	}

	resp.Id = user.ID
	return resp, nil
}
