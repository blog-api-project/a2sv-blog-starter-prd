package routers

import (
	"blog_api/Delivery/controllers"
	contracts_services "blog_api/Domain/contracts/services"
	infrastructure "blog_api/Infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	userController *controllers.UserController,
	tokenController *controllers.TokenController,
	oauthController *controllers.OAuthController,
	adminController *controllers.AdminController,
	blogController *controllers.BlogController,
	commentController *controllers.CommentController,
	aiController *controllers.AIController, // Added AI controller
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

	// OAuth routes
	oauthRoutes := router.Group("/api/auth")
	{
		oauthRoutes.GET("/:provider/login", oauthController.InitiateOAuthFlow)
		oauthRoutes.GET("/:provider/callback", oauthController.HandleOAuthCallback)
		oauthRoutes.POST("/:provider/link", 
			infrastructure.AuthMiddleware(jwtService),
			infrastructure.RBACMiddleware("user", "admin"), 
			oauthController.LinkOAuthToExistingUser)
	}

	// Admin routes
	adminRoutes := router.Group("/api/admin")
	adminRoutes.Use(
		infrastructure.AuthMiddleware(jwtService), 
		infrastructure.RBACMiddleware("admin"),
	)
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
		blogRoutes.GET("/search", blogController.SearchBlogsHandler)
		blogRoutes.POST("/:id/like", blogController.LikeBlog)
		blogRoutes.POST("/:id/dislike", blogController.DislikeBlog)
		blogRoutes.POST("/:id/generate-content",
			infrastructure.RBACMiddleware("user", "admin"),
			aiController.GenerateBlogContentForPost)
	}

	// Comment routes
	commentRoutes := router.Group("/api/comments")
	commentRoutes.Use(infrastructure.AuthMiddleware(jwtService))
	{ 
		commentRoutes.POST("/create/:id", commentController.CreateComment)
		commentRoutes.PUT("/:id", commentController.UpdateComment)
		commentRoutes.DELETE("/:id", commentController.DeleteComment)
	}

	// AI routes
	aiRoutes := router.Group("/api/ai")
	aiRoutes.Use(infrastructure.AuthMiddleware(jwtService))
	{
		aiRoutes.POST("/generate", 
			infrastructure.RBACMiddleware("user", "admin"),
			aiController.GenerateBlogPost)
		aiRoutes.POST("/suggest-improvements", 
			infrastructure.RBACMiddleware("user", "admin"),
			aiController.SuggestImprovements)
	}

	return router
}
