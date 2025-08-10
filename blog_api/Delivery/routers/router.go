package routers

import (
	"blog_api/Delivery/controllers"
	contracts_services "blog_api/Domain/contracts/services"
	infrastructure "blog_api/Infrastructure"

	"github.com/gin-gonic/gin"
)

// SetupRouter sets up the application routes
func SetupRouter(
	userController *controllers.UserController,
	tokenController *controllers.TokenController,
	adminController *controllers.AdminController,
	blogController *controllers.BlogController,
	jwtService contracts_services.IJWTService,
) *gin.Engine {
	router := gin.Default()

	// User routes
	userRoutes := router.Group("/api/users")
	{
		userRoutes.POST("/register", userController.Register)
		userRoutes.POST("/login", userController.Login)
		userRoutes.POST("/logout", userController.Logout)
		userRoutes.POST("/forgot-password", userController.ForgotPassword)
		userRoutes.POST("/reset-password", userController.ResetPassword)

		// Auth required
		userRoutes.Use(infrastructure.AuthMiddleware(jwtService))
		userRoutes.PUT("/profile", userController.UpdateProfile)
	}

	// Authentication routes
	authRoutes := router.Group("/api/auth")
	{
		authRoutes.POST("/refresh", tokenController.RefreshToken)
		authRoutes.POST("/validate", tokenController.ValidateToken)
	}

	
	// Admin routes
	adminRoutes := router.Group("/api/admin")
	adminRoutes.Use(infrastructure.AuthMiddleware(jwtService), infrastructure.RBACMiddleware("admin"))
	{
		adminRoutes.POST("/users/:userID/promote", adminController.PromoteUser)
		adminRoutes.POST("/users/:userID/demote", adminController.DemoteUser)
	}

	// Blog routes
	blogRoutes := router.Group("/api/blogs")
	blogRoutes.Use(infrastructure.AuthMiddleware(jwtService))
	{
		blogRoutes.POST("/create", blogController.CreateBlog)
		blogRoutes.GET("/", blogController.GetBlogs)
		blogRoutes.PUT("/:id", blogController.UpdateBlogHandler)
		blogRoutes.DELETE("/:id", blogController.DeleteBlogHandler)
	}

	return router
}
