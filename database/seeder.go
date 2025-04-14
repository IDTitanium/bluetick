package database

import (
	"bluetick/models"
)

// Seed data
func Seed() {
	// Check if user already exists
	var count int64
	DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		return
	}

	user := models.User{Username: "user1"}
	user.SetPassword("password1")
	DB.Create(&user)
}
