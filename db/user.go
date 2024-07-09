package db

import (
	"context"
	"strconv"
	"time"

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
	SpotifyToken        string    // This is the access token, we will use it to request user-specific data
	SpotifyRefreshToken string    // This is the refresh token, we will use it to get a new access token
	SpotifyLastUpdated  time.Time // This is the last time we updated the user's Spotify data
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
		SpotifyLastUpdated:  time.Now().Add(-time.Hour * 48), // We want to queue the user for an update soon, so we set this to 48 hours ago
	}
	if err := DB.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUpdateableUsers gets users that need to be updated
func GetUpdateableUsers(ctx context.Context) ([]*User, error) {
	// Retrieve the relevant keys from the config service
	s_interval := clients.GetKey("CRAWLER_INTERVAL_USER")
	interval, err := strconv.ParseInt(s_interval, 10, 64)
	if err != nil {
		return nil, err
	}

	s_limit := clients.GetKey("CRAWLER_BATCH_SIZE")
	limit, err := strconv.ParseInt(s_limit, 10, 64)
	if err != nil {
		return nil, err
	}

	// Calculate the time that we need to look back to
	lookback := time.Now().Add(-time.Second * time.Duration(interval))

	// Get the users that need to be updated
	var users []*User
	if err := DB.Where("spotify_last_updated < ?", lookback).Limit(int(limit)).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
