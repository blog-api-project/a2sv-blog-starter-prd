package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Returns context with default timeout for database operations
func DefaultTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func ConnectToMongoDB() (*mongo.Database, error) {
	// Get MongoDB URI from environment variable or use default
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	// Get database name from environment variable or use default
	dbName := os.Getenv("MONGODB_DB_NAME")
	if dbName == "" {
		dbName = "blog_platform"
	}
	clientOptions := options.Client().ApplyURI(mongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to MongoDB!")

	return client.Database(dbName), nil
}

//Closes the MongoDB connection
func CloseMongoDBConnection(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return client.Disconnect(ctx)
} 