package routes

import (
	"bluetick/controllers"
	"bluetick/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterMessageRoutes(router *gin.Engine) {
	group := router.Group("/messages")
	group.Use(middlewares.AuthMiddleware())

	group.GET("/", controllers.GetMessages)
	group.PUT("/read", controllers.MarkMessagesAsRead)
}
