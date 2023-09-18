package server

import (
	"context"

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
		return resp, errNotFound
	}

	if s.exits(ctx, username) {
		return resp, errAlreadyExist
	}

	user := &models.User{
		ID:       id,
		Email:    email,
		Username: username,
		Password: []byte(password),
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

	return resp, nil
}
