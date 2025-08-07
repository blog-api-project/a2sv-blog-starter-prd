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

  //Blog Router
  blogRoutes := router.Group("/api/blogs")
  blogRoutes.Use(infrastructure.AuthMiddleware(jwtService))

  {
         blogRoutes.POST("/create", blogController.CreateBlog)
		 blogRoutes.GET("/",blogController.GetBlogs)
		 blogRoutes.PUT("/:id",blogController.UpdateBlogHandler)
	     blogRoutes.DELETE("/:id",blogController.DeleteBlogHandler)
  }

  return router
}
