package controllers

import (
	"bluetick/database"
	"bluetick/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetMessages retrieves chat history between auth user and another user
func GetMessages(c *gin.Context) {
	userID := c.GetUint64("user_id")
	otherID, err := strconv.Atoi(c.Query("other_id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid other_id"})
		return
	}

	var messages []models.Message

	database.DB.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		userID, otherID, otherID, userID).Order("created_at asc").Find(&messages)

	c.JSON(http.StatusOK, messages)
}
