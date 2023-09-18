package server

import (
	"context"
	"encoding/base64"
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

	if !s.valid(md["authorization"], info.FullMethod) {
		return nil, errInvalidToken
	}
	return handler(ctx, req)
}

func (s *Server) valid(authorization []string, method string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Basic ")

	basic, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return false
	}

	auth := strings.Split(string(basic), ":")
	if len(auth) < 2 {
		return false
	}
	username := auth[0]
	password := auth[1]

	if !s.Exist(username, password) {
		return false
	}

	switch method {
	case "/authgrpc.Users/Create", "/authgrpc.Users/Update", "/authgrpc.Users/Delete":
		if !s.Admin(username) {
			return false
		}
	}

	return true
}

// It's not very efficient, I know
func (s *Server) Exist(username, password string) bool {
	ids, err := s.db.List(context.Background())
	if err != nil {
		return false
	}

	for _, id := range ids {
		entry, err := s.db.Get(context.Background(), id)
		if err != nil {
			return false
		}

		user, ok := entry.Val.(*models.User)
		if !ok {
			return false
		}

		if user.Username == username && string(user.Password) == password {
			return true
		}
	}
	return false
}

func (s *Server) Admin(username string) bool {
	ids, err := s.db.List(context.Background())
	if err != nil {
		return false
	}

	for _, id := range ids {
		entry, err := s.db.Get(context.Background(), id)
		if err != nil {
			return false
		}

		user, ok := entry.Val.(*models.User)
		if !ok {
			return false
		}

		if user.Username == username {
			return user.Admin
		}
	}
	return false
}
