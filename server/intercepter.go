package server

import (
	"context"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/valli0x/auth-grpc/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

func (s *Server) ValidToken(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}

	if !s.valid(ctx, md["authorization"], info.FullMethod) {
		return nil, errInvalidToken
	}
	return handler(ctx, req)
}

func (s *Server) valid(ctx context.Context, authorization []string, method string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Basic ")

	basic, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return false
	}

	auth := strings.Split(string(basic), ":")
	if len(auth) != 2 {
		return false
	}
	username := auth[0]
	password := auth[1]

	if !s.checkPass(ctx, username, password) {
		return false
	}

	switch method {
	case "/authgrpc.Users/Create", "/authgrpc.Users/Update", "/authgrpc.Users/Delete":
		if !s.admin(ctx, username) {
			return false
		}
	}

	return true
}

func (s *Server) checkPass(ctx context.Context, username, password string) bool {
	user, err := s.userByUsername(ctx, username)
	if err != nil {
		return false
	}
	if user.Username == username && string(user.Password) == password {
		return true
	}
	return false
}

func (s *Server) exits(ctx context.Context, username string) bool {
	user, err := s.userByUsername(ctx, username)
	if err != nil {
		return false
	}
	if user.Username == username {
		return true
	}
	return false
}

func (s *Server) admin(ctx context.Context, username string) bool {
	user, err := s.userByUsername(ctx, username)
	if err != nil {
		return false
	}
	if user.Username == username {
		return user.Admin
	}
	return false
}

func (s *Server) userByUsername(ctx context.Context, username string) (*models.User, error) {
	entry, err := s.db.Get(ctx, username)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, errNotFound
	}
	user, ok := entry.Val.(*models.User)
	if !ok {
		return nil, errors.New("converting error")
	}
	return user, nil
}
