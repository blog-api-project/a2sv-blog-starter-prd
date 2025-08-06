package repositories

import (
	repositories "blog_api/Domain/contracts/repositories"
	"blog_api/Domain/models"
	"blog_api/Repositories/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type mongoTokenRepository struct {
	accessTokensCollection  *mongo.Collection
	refreshTokensCollection *mongo.Collection
}


func NewMongoTokenRepository(accessCollection, refreshCollection *mongo.Collection) repositories.ITokenRepository {
	return &mongoTokenRepository{
		accessTokensCollection:  accessCollection,
		refreshTokensCollection: refreshCollection,
	}
}

// stores an access token in the database
func (r *mongoTokenRepository) StoreAccessToken(accessToken *models.AccessToken) error {
	ctx, cancel :=database. DefaultTimeout()
	defer cancel()

	_, err := r.accessTokensCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{"token": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		return err
	}
	objectID := primitive.NewObjectID()
	accessToken.ID = objectID.Hex()

	doc := bson.M{
		"_id":        objectID,
		"user_id":    accessToken.UserID,
		"token":      accessToken.Token,
		"expires_at": accessToken.ExpiresAt,
		"created_at": accessToken.CreatedAt,
		"updated_at": accessToken.UpdatedAt,
	}

	_, err = r.accessTokensCollection.InsertOne(ctx, doc)
	return err
}

//stores a refresh token in the database
func (r *mongoTokenRepository) StoreRefreshToken(refreshToken *models.RefreshToken) error {
	ctx, cancel :=database. DefaultTimeout()
	defer cancel()

	// Create index on token field for faster lookups
	_, err := r.refreshTokensCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{"token": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		return err
	}

	// Generate ObjectID for the token
	objectID := primitive.NewObjectID()
	refreshToken.ID = objectID.Hex()

	doc := bson.M{
		"_id":        objectID,
		"user_id":    refreshToken.UserID,
		"token":      refreshToken.Token,
		"expires_at": refreshToken.ExpiresAt,
		"created_at": refreshToken.CreatedAt,
		"updated_at": refreshToken.UpdatedAt,
	}

	_, err = r.refreshTokensCollection.InsertOne(ctx, doc)
	return err
}

//validates if an access token exists and is not expired
func (r *mongoTokenRepository) ValidateAccessToken(token string) (bool, error) {
	ctx, cancel := database.DefaultTimeout()
	defer cancel()

	var doc bson.M
	err := r.accessTokensCollection.FindOne(ctx, bson.M{
		"token": token,
		"expires_at": bson.M{
			"$gt": time.Now(),
		},
	}).Decode(&doc)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

//validates if a refresh token exists and is not expired
func (r *mongoTokenRepository) ValidateRefreshToken(token string) (bool, error) {
	ctx, cancel := database.DefaultTimeout()
	defer cancel()

	var doc bson.M
	err := r.refreshTokensCollection.FindOne(ctx, bson.M{
		"token": token,
		"expires_at": bson.M{
			"$gt": time.Now(),
		},
	}).Decode(&doc)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

//revokes an access token
func (r *mongoTokenRepository) RevokeAccessToken(token string) error {
	ctx, cancel :=database. DefaultTimeout()
	defer cancel()

	_, err := r.accessTokensCollection.DeleteOne(ctx, bson.M{"token": token})
	return err
}

//revokes a refresh token
func (r *mongoTokenRepository) RevokeRefreshToken(token string) error {
	ctx, cancel :=database. DefaultTimeout()
	defer cancel()

	_, err := r.refreshTokensCollection.DeleteOne(ctx, bson.M{"token": token})
	return err
}

//revokes all tokens for a user
func (r *mongoTokenRepository) RevokeAllUserTokens(userID string) error {
	ctx, cancel :=database. DefaultTimeout()
	defer cancel()
	_, err := r.accessTokensCollection.DeleteMany(ctx, bson.M{"user_id": userID})
	if err != nil {
		return err
	}
	_, err = r.refreshTokensCollection.DeleteMany(ctx, bson.M{"user_id": userID})
	return err
} 