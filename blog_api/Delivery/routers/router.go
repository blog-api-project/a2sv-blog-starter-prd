package routers

import (
  "blog_api/Delivery/controllers"
  "blog_api/Domain/contracts/services"
  infrastructure "blog_api/Infrastructure"

  "github.com/gin-gonic/gin"
)

func SetupRouter(userController *controllers.UserController, tokenController *controllers.TokenController,blogController *controllers.BlogController,jwtService services.IJWTService,
) *gin.Engine {
  router := gin.Default()

<<<<<<< HEAD
  // User routes
  userRoutes := router.Group("/api/users")
  {
    userRoutes.POST("/register", userController.Register)
    userRoutes.POST("/login", userController.Login)
    userRoutes.POST("/logout", userController.Logout)
    userRoutes.POST("/forgot-password", userController.ForgotPassword)
    userRoutes.POST("/reset-password", userController.ResetPassword)
  }
=======
	// User routes
	userRoutes := router.Group("/api/users")
	{
		userRoutes.POST("/register", userController.Register)
		userRoutes.POST("/login", userController.Login)
		userRoutes.POST("/logout", userController.Logout)
		userRoutes.POST("/forgot-password", userController.ForgotPassword)
		userRoutes.POST("/reset-password", userController.ResetPassword)
	}
>>>>>>> 3e802a94ec285b5614213e9bb3ea3fe693d3ebad

  // Authentication routes
  authRoutes := router.Group("/api/auth")
  {
    authRoutes.POST("/refresh", tokenController.RefreshToken)
    authRoutes.POST("/validate", tokenController.ValidateToken)
  }

  //Blog Router
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
