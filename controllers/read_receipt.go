package controllers

import (
	"bluetick/database"
	"bluetick/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MarkMessagesAsRead updates all messages as read for a chat
func MarkMessagesAsRead(c *gin.Context) {
	senderID, _ := strconv.Atoi(c.Query("sender_id"))
	receiverID := c.GetUint64("user_id")

	database.DB.Model(&models.Message{}).
		Where("sender_id = ? AND receiver_id = ? AND is_read = false", senderID, receiverID).
		Update("is_read", true)

	c.JSON(http.StatusOK, gin.H{"message": "Messages marked as read"})
}
