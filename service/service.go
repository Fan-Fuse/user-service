package service

import (
	"context"
	"strconv"

	"github.com/Fan-Fuse/user-service/db"
	"github.com/Fan-Fuse/user-service/proto"
	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedUserServiceServer
}

// RegisterServer registers the server with the gRPC server
func RegisterServer(s *grpc.Server) {
	proto.RegisterUserServiceServer(s, &server{})
}

func (s *server) GetUser(ctx context.Context, in *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	user, err := db.GetUser(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &proto.GetUserResponse{
		Id:       strconv.Itoa(int(user.ID)),
		Name:     user.Name,
		ImageUrl: user.ImageURL,
	}, nil
}

func (s *server) CreateUser(ctx context.Context, in *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	user, err := db.CreateUser(ctx, in.Name, in.ImageUrl, in.SpotifyUser.Id, in.SpotifyUser.AccessToken, in.SpotifyUser.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &proto.CreateUserResponse{
		Success: true,
		Id:      strconv.Itoa(int(user.ID)),
	}, nil
}
