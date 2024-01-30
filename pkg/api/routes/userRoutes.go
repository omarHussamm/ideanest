package routes

import (
	"demo/pkg/api/handlers"

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	router.POST("/signup", handlers.SignUp)
	router.POST("/signin", handlers.SignIn)
	router.POST("/refresh-token", handlers.RefreshToken)
	router.GET("/users", handlers.GetAllUsers)
}
