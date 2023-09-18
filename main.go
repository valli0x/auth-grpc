package main

import (
	"context"
	"log"
	"net"

	"github.com/hashicorp/go-uuid"
	"github.com/valli0x/auth-grpc/models"
	"github.com/valli0x/auth-grpc/server"
	"github.com/valli0x/auth-grpc/storage"
	"github.com/valli0x/auth-grpc/storage/inmem"
	pb "github.com/valli0x/grpc-proto/auth-grpc/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port          = ":8080"
	adminUsername = "admin"
	adminPass     = "admin"
)

func main() {
	// create database
	db, err := inmem.NewInmem()
	if err != nil {
		panic(err)
	}
	// create admin user
	err = rootUser(db)
	if err != nil {
		panic(err)
	}
	// create users grpc handler
	s, err := server.NewServer(db)
	if err != nil {
		panic(err)
	}
	// options grpc server
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(s.ValidToken),
	}
	// start grpc server
	grpcServer := grpc.NewServer(opts...)
	reflection.Register(grpcServer)
	pb.RegisterUsersServer(grpcServer, s)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func rootUser(db storage.Storage) error {
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}

	user := &models.User{
		ID:       uuid,
		Email:    "hello@bel.by",
		Username: adminUsername,
		Password: []byte(adminPass),
		Admin:    true,
	}

	if err := db.Put(context.Background(), &storage.Entry{
		Key: user.ID,
		Val: user,
	}); err != nil {
		return err
	}

	return nil
}
