package db

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name                string
	Age                 int
	ImageURL            string
	SpotifyID           string
	SpotifyToken        string // This is the access token, we will use it to request user-specific data
	SpotifyRefreshToken string // This is the refresh token, we will use it to get a new access token
}
