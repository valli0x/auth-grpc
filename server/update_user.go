package server

import (
	"context"
	"errors"

	"github.com/valli0x/auth-grpc/models"
	"github.com/valli0x/auth-grpc/storage"
	pb "github.com/valli0x/grpc-proto/auth-grpc/api"
)

func (s *Server) Update(ctx context.Context, r *pb.UpdateRequest) (*pb.Empty, error) {
	resp := &pb.Empty{}

	id := r.User.GetId()
	email := r.User.GetEmail()
	username := r.User.GetUsername()
	password := r.User.GetPassword()

	entry, err := s.db.Get(ctx, id)
	if err != nil {
		return resp, err
	}

	if entry == nil {
		return resp, errors.New(UserNotFoundError)
	}

	if s.exits(username) {
		return resp, errors.New("username is not unique")
	}

	user := &models.User{
		ID:       id,
		Email:    email,
		Username: username,
		Password: []byte(password),
	}

	if err := s.db.Put(ctx, &storage.Entry{
		Key: user.ID,
		Val: user,
	}); err != nil {
		return resp, err
	}

	return resp, nil
}
