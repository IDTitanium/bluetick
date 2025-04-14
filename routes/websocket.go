package routes

import (
	"bluetick/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterWebSocketRoutes(router *gin.Engine) {
	router.GET("/ws", controllers.HandleWebSocket)
}
