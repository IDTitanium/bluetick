package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserStatus returns online status of a user
func GetUserStatus(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("user_id"))

	mu.Lock()
	online := usersOnline[uint64(userID)]
	mu.Unlock()

	c.JSON(http.StatusOK, gin.H{"user_id": userID, "online": online})
}
