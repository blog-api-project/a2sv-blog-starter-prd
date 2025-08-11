package repositories

import (
	"blog_api/Domain/contracts/repositories"
	"blog_api/Domain/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoOAuthRepository implements IOAuthRepository for MongoDB
type MongoOAuthRepository struct {
	collection *mongo.Collection
}

// NewMongoOAuthRepository creates a new MongoDB OAuth repository
func NewMongoOAuthRepository(collection *mongo.Collection) repositories.IOAuthRepository {
	return &MongoOAuthRepository{
		collection: collection,
	}
}

// CreateOAuthUser creates a new OAuth user record
func (r *MongoOAuthRepository) CreateOAuthUser(oauthUser *models.OAuthUser) error {
	oauthUser.ID = primitive.NewObjectID().Hex()
	oauthUser.CreatedAt = time.Now()
	oauthUser.UpdatedAt = time.Now()
	
	_, err := r.collection.InsertOne(context.TODO(), oauthUser)
	return err
}

// GetOAuthUserByProviderID retrieves OAuth user by provider ID
func (r *MongoOAuthRepository) GetOAuthUserByProviderID(provider, providerID string) (*models.OAuthUser, error) {
	filter := bson.M{
		"provider":     provider,
		"provider_id": providerID,
	}
	
    var oauthUser models.OAuthUser
    err := r.collection.FindOne(context.TODO(), filter).Decode(&oauthUser)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, mongo.ErrNoDocuments
        }
        return nil, err
    }
	
	return &oauthUser, nil
}

// GetOAuthUserByEmail retrieves OAuth user by email and provider
func (r *MongoOAuthRepository) GetOAuthUserByEmail(provider, email string) (*models.OAuthUser, error) {
	filter := bson.M{
		"provider": provider,
		"email":    email,
	}
	
    var oauthUser models.OAuthUser
    err := r.collection.FindOne(context.TODO(), filter).Decode(&oauthUser)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, mongo.ErrNoDocuments
        }
        return nil, err
    }
	
	return &oauthUser, nil
}

// UpdateOAuthUser updates an existing OAuth user record
func (r *MongoOAuthRepository) UpdateOAuthUser(oauthUser *models.OAuthUser) error {
	oauthUser.UpdatedAt = time.Now()
	
	filter := bson.M{"_id": oauthUser.ID}
	update := bson.M{"$set": oauthUser}
	
	_, err := r.collection.UpdateOne(context.TODO(), filter, update)
	return err
}

// LinkOAuthToUser links an OAuth user to an existing user account
func (r *MongoOAuthRepository) LinkOAuthToUser(oauthUserID, userID string) error {
	filter := bson.M{"_id": oauthUserID}
	update := bson.M{
		"$set": bson.M{
			"user_id":   userID,
			"updated_at": time.Now(),
		},
	}
	
	_, err := r.collection.UpdateOne(context.TODO(), filter, update)
	return err
} 