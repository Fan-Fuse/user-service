package db

import (
	"context"

	"github.com/Fan-Fuse/user-service/clients"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name                string
	Age                 int
	ImageURL            string
	SpotifyID           string
	SpotifyToken        string // This is the access token, we will use it to request user-specific data
	SpotifyRefreshToken string // This is the refresh token, we will use it to get a new access token
}

// GetUser gets a user by ID
func GetUser(ctx context.Context, id string) (*User, error) {
	var user User
	if err := DB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a user
func CreateUser(ctx context.Context, name, imageURL, spotifyID, spotifyToken, spotifyRefreshToken string) (*User, error) {
	// Check if we allow users to register
	if clients.GetKey("USER_REGISTRATION_OPEN") != "1" {
		return nil, status.Error(codes.PermissionDenied, "user registration is not open")
	}

	// Make sure the user does not already exist
	var count int64
	if err := DB.Model(&User{}).Where("spotify_id = ?", spotifyID).Count(&count).Error; err != nil {
		return nil, status.Error(codes.Internal, "could not check if user exists")
	}
	if count > 0 {
		// Retrieve the user
		var user User
		if err := DB.Where("spotify_id = ?", spotifyID).First(&user).Error; err != nil {
			return nil, status.Error(codes.Internal, "could not retrieve user")
		}
		return &user, nil
	}

	user := User{
		Name:                name,
		ImageURL:            imageURL,
		SpotifyID:           spotifyID,
		SpotifyToken:        spotifyToken,
		SpotifyRefreshToken: spotifyRefreshToken,
	}
	if err := DB.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
