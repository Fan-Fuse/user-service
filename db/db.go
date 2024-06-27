package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name                string
	Age                 int
	SpotifyID           string
	SpotifyToken        string // This is the access token, we will use it to request user-specific data
	SpotifyRefreshToken string // This is the refresh token, we will use it to get a new access token
}

var DB *gorm.DB

func Init() {
	db, err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&User{})

	DB = db
}
