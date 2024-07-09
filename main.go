package main

import (
	"net"
	"os"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/Fan-Fuse/user-service/clients"
	"github.com/Fan-Fuse/user-service/db"
	pb "github.com/Fan-Fuse/user-service/proto"
)

type server struct {
	pb.UnimplementedUserServiceServer
}

func init() {
	// Initialize logger
	logger := zap.Must(zap.NewProduction())
	if os.Getenv("APP_ENV") == "development" {
		logger = zap.Must(zap.NewDevelopment())
	}

	zap.ReplaceGlobals(logger)

	// Initialize config service client
	clients.InitConfig(os.Getenv("CONFIG_ADDRESS"))

	// Initialize database
	db.Init()
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		zap.S().Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{})

	zap.S().Info("Server started on port 50051")
	if err := s.Serve(lis); err != nil {
		zap.S().Fatalf("Failed to serve: %v", err)
	}
}
