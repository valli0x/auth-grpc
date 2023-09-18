package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-uuid"
	"github.com/spf13/viper"
	"github.com/valli0x/auth-grpc/models"
	"github.com/valli0x/auth-grpc/server"
	"github.com/valli0x/auth-grpc/storage"
	"github.com/valli0x/auth-grpc/storage/inmem"
)

const (
	address       = ":8080"
	adminEmail    = "admin@minsk.by"
	adminUsername = "admin"
	adminPass     = "admin"
)

func main() {
	viper.BindEnv("port", "GRPC_PORT")
	viper.SetDefault("port", address)
	port := viper.GetString("port")

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
	fmt.Printf("grpc server starting on port %s...\n", port)
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
		Email:    adminEmail,
		Username: adminUsername,
		Password: []byte(adminPass),
		Admin:    true,
	}

	// by ID
	if err := db.Put(context.Background(), &storage.Entry{
		Key: user.ID,
		Val: user,
	}); err != nil {
		return err
	}

	// by Username
	if err := db.Put(context.Background(), &storage.Entry{
		Key: user.Username,
		Val: user,
	}); err != nil {
		return err
	}

	return nil
}
