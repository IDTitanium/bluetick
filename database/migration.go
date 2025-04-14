package database

import (
	"bluetick/models"
	"fmt"
)

func MigrateDB() {
	err := DB.AutoMigrate(&models.User{}, &models.Message{})
	if err != nil {
		fmt.Println("Database migration failed:", err)
	} else {
		fmt.Println("Database migrated successfully!")
	}
}
