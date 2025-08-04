package routers

import (
	"blog_api/Delivery/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userController *controllers.UserController, tokenController *controllers.TokenController) *gin.Engine {
	router := gin.Default()

	// User routes
	userRoutes := router.Group("/api/users")
	{
		userRoutes.POST("/register", userController.Register)
		userRoutes.POST("/login", userController.Login)
		userRoutes.POST("/logout", userController.Logout)
	}

	// Authentication routes
	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/refresh", tokenController.RefreshToken)
		authRoutes.POST("/validate", tokenController.ValidateToken)
	}

	return router
} 