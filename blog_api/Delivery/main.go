package main

import (
	"blog_api/Delivery/controllers"
	"blog_api/Delivery/routers"
	infrastructure "blog_api/Infrastructure"
	repositories "blog_api/Repositories"
	"blog_api/Repositories/database"
	usecases "blog_api/Usecases"

	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load environment variables
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
	roleRepo := repositories.NewMongoRoleRepository(db.Collection("roles"))
	blogRepo := repositories.NewMongoBlogRepository(db.Collection("Blogs"))

	// Initialize services
	passwordSvc := infrastructure.NewPasswordService()
	jwtSvc := infrastructure.NewJWTService()
	validationSvc := infrastructure.NewValidationService()
	emailSvc := infrastructure.NewEmailService()
	imageSvc := infrastructure.NewImageService(uploadDir)

	// Dev seeding: roles + initial admin account (non-production only)
	if os.Getenv("ENV") != "production" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		rolesCol := db.Collection("roles")
		usersCol := db.Collection("users")

		// Ensure roles exist (user, admin)
		_, _ = rolesCol.UpdateOne(
			ctx,
			bson.M{"role": "user"},
			bson.M{"$setOnInsert": bson.M{"role": "user", "created_at": time.Now(), "updated_at": time.Now()}},
			options.Update().SetUpsert(true),
		)
		_, _ = rolesCol.UpdateOne(
			ctx,
			bson.M{"role": "admin"},
			bson.M{"$setOnInsert": bson.M{"role": "admin", "created_at": time.Now(), "updated_at": time.Now()}},
			options.Update().SetUpsert(true),
		)
		var adminRoleDoc bson.M
		if err := rolesCol.FindOne(ctx, bson.M{"role": "admin"}).Decode(&adminRoleDoc); err == nil {
			if adminOID, ok := adminRoleDoc["_id"].(primitive.ObjectID); ok {
				// Ensure first admin user exists (admin@example.com)
				hashed, _ := passwordSvc.HashPassword("Admin#12345")
				set := bson.M{
					"username":       "admin_user",
					"first_name":     "Admin",
					"last_name":      "One",
					"password":       hashed,
					"role_id":        adminOID,
					"is_active":      true,
					"email_verified": true,
					"updated_at":     time.Now(),
				}
				update := bson.M{
					"$set":         set,
					"$setOnInsert": bson.M{"created_at": time.Now()},
				}
				_, _ = usersCol.UpdateOne(ctx, bson.M{"email": "admin@example.com"}, update, options.Update().SetUpsert(true))
			}
		}
	}

	

	// Initialize use cases
	tokenUseCase := usecases.NewTokenUseCase(tokenRepo, jwtSvc, roleRepo)
	userUseCase := usecases.NewUserUseCase(userRepo, passwordSvc, jwtSvc, validationSvc, emailSvc, tokenUseCase, roleRepo)
	adminUseCase := usecases.NewAdminUseCase(userRepo, roleRepo)
	blogUseCase := usecases.NewBlogUseCase(blogRepo)

	// Initialize controllers
	userController := controllers.NewUserController(userUseCase, tokenUseCase, jwtSvc)
	tokenController := controllers.NewTokenController(tokenUseCase, jwtSvc)
	adminController := controllers.NewAdminController(adminUseCase)
	blogController := controllers.NewBlogController(blogUseCase, imageSvc)

	// Setup router
	router := routers.SetupRouter(userController, tokenController, adminController, blogController, jwtSvc)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
