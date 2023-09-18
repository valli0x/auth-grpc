package main

import (
	"context"
	"log"

	"github.com/hashicorp/go-uuid"
	"github.com/valli0x/auth-grpc/models"
	"github.com/valli0x/auth-grpc/server"
	"github.com/valli0x/auth-grpc/storage"
	"github.com/valli0x/auth-grpc/storage/inmem"
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
		log.Fatalf("failed to create storage: %v", err)
	}

	// create admin user
	err = rootUser(db)
	if err != nil {
		log.Fatalf("failed to create root user: %v", err)
	}

	// create grpc server
	grpcServer, err := server.NewServer(db)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}
	
	// start grpc server
	if err := grpcServer.RunServer(port); err != nil {
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
