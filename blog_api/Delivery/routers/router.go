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
		userRoutes.POST("/forgot-password", userController.ForgotPassword)
		userRoutes.POST("/reset-password", userController.ResetPassword)
	}

	// Authentication routes
	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/refresh", tokenController.RefreshToken)
		authRoutes.POST("/validate", tokenController.ValidateToken)
	}
	// Admin routes (requires auth + RBAC admin)
	adminController := controllers.NewAdminController(adminUseCase)
	adminRoutes := router.Group("/api/admin")
    adminRoutes.Use(infrastructure.AuthMiddleware(), infrastructure.RBACMiddleware("admin"))
	{
		adminRoutes.POST("/users/:userID/promote", adminController.PromoteUser)
		adminRoutes.POST("/users/:userID/demote", adminController.DemoteUser)
	}



	return router
} 