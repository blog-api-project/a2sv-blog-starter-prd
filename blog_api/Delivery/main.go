package main

import (
	"blog_api/Delivery/controllers"
	"blog_api/Delivery/routers"
	infrastructure "blog_api/Infrastructure"
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

	// Initialize repositories
	userRepo := database.NewMongoUserRepository(db.Collection("users"))
	tokenRepo := database.NewMongoTokenRepository(db.Collection("access_tokens"), db.Collection("refresh_tokens"))
	
	// Initialize services
	passwordSvc := infrastructure.NewPasswordService()
	jwtSvc := infrastructure.NewJWTService()
	validationSvc := infrastructure.NewValidationService()
	
	// Initialize use cases
	tokenUseCase := usecases.NewTokenUseCase(tokenRepo, jwtSvc)
	userUseCase := usecases.NewUserUseCase(userRepo, passwordSvc, jwtSvc, validationSvc, tokenUseCase)
	
	// Initialize controllers
	userController := controllers.NewUserController(userUseCase, tokenUseCase, jwtSvc)
	tokenController := controllers.NewTokenController(tokenUseCase, jwtSvc)
	
	router := routers.SetupRouter(userController, tokenController)
	
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