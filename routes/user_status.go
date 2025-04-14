package routes

import (
	"bluetick/controllers"
	"bluetick/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterUserStatusRoutes(router *gin.Engine) {
	group := router.Group("/user")
	group.Use(middlewares.AuthMiddleware())
	group.GET("/user/status/:user_id", controllers.GetUserStatus)
}
