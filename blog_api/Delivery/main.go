package main

import (
  "blog_api/Delivery/controllers"
  "blog_api/Delivery/routers"
  infrastructure "blog_api/Infrastructure"
  repositories "blog_api/Repositories"
  "blog_api/Repositories/database"
  usecases "blog_api/Usecases"
  "log"
  "os"

  "github.com/joho/godotenv"
)

func main() {
  if err := godotenv.Load(); err != nil {
    log.Printf("Warning: Error loading .env file: %v", err)
  }

  // Connect to MongoDB
  db, err := database.ConnectToMongoDB()
  if err != nil {
    log.Fatal("Failed to connect to MongoDB:", err)
  }
  defer func() {
    if client := db.Client(); client != nil {
      if err := database.CloseMongoDBConnection(client); err != nil {
        log.Printf("Error closing MongoDB connection: %v", err)
      }
    }
  }()
  // Setup upload directory
  uploadDir := os.Getenv("IMAGE_UPLOAD_PATH")
  if uploadDir == "" {
    uploadDir = "../uploads"
  }
  


  // Initialize repositories
  userRepo := repositories.NewMongoUserRepository(db.Collection("users"))
  tokenRepo := repositories.NewMongoTokenRepository(db.Collection("access_tokens"), db.Collection("refresh_tokens"))
  blogRepo := repositories.NewMongoBlogRepository(db.Collection("Blogs"))
  
  // Initialize services
  passwordSvc := infrastructure.NewPasswordService()
  jwtSvc := infrastructure.NewJWTService()
  validationSvc := infrastructure.NewValidationService()
  imageSvc := infrastructure.NewImageService(uploadDir)
  
  // Initialize use cases
  tokenUseCase := usecases.NewTokenUseCase(tokenRepo, jwtSvc)
  userUseCase := usecases.NewUserUseCase(userRepo, passwordSvc, jwtSvc, validationSvc, tokenUseCase)
  blogUseCase := usecases.NewBlogUseCase(blogRepo)
  
  // Initialize controllers
  userController := controllers.NewUserController(userUseCase, tokenUseCase, jwtSvc)
  tokenController := controllers.NewTokenController(tokenUseCase, jwtSvc)
  blogController := controllers.NewBlogController(blogUseCase,imageSvc)
  
  router := routers.SetupRouter(userController, tokenController, blogController, jwtSvc)
  
  // Get port from environment variable or use default
  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
  }
  log.Printf("Server starting on port %s...", port)
  if err := router.Run(":" + port); err != nil {
    log.Fatal("Failed to start server:", err)
  }
}
