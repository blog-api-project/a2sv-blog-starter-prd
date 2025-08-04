package database

import (
	repositories "blog_api/Domain/contracts/repositories"
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

//Stores an access token in the database
func (r *mongoTokenRepository) StoreAccessToken(userID, token string, expiresAt time.Time) error {
	ctx, cancel := DefaultTimeout()
	defer cancel()

	_, err := r.accessTokensCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{"token": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		return err
	}

	doc := bson.M{
		"_id":        primitive.NewObjectID(),
		"user_id":    userID,
		"token":      token,
		"expires_at": expiresAt,
		"created_at": time.Now(),
	}

	_, err = r.accessTokensCollection.InsertOne(ctx, doc)
	return err
}

//Stores a refresh token in the database
func (r *mongoTokenRepository) StoreRefreshToken(userID, token string, expiresAt time.Time) error {
	ctx, cancel := DefaultTimeout()
	defer cancel()

	_, err := r.refreshTokensCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{"token": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		return err
	}

	doc := bson.M{
		"_id":        primitive.NewObjectID(),
		"user_id":    userID,
		"token":      token,
		"expires_at": expiresAt,
		"created_at": time.Now(),
	}

	_, err = r.refreshTokensCollection.InsertOne(ctx, doc)
	return err
}

//Validates if an access token exists and is not expired
func (r *mongoTokenRepository) ValidateAccessToken(token string) (bool, error) {
	ctx, cancel := DefaultTimeout()
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

//Validates if a refresh token exists and is not expired
func (r *mongoTokenRepository) ValidateRefreshToken(token string) (bool, error) {
	ctx, cancel := DefaultTimeout()
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
	ctx, cancel := DefaultTimeout()
	defer cancel()

	_, err := r.accessTokensCollection.DeleteOne(ctx, bson.M{"token": token})
	return err
}

//revokes a refresh token
func (r *mongoTokenRepository) RevokeRefreshToken(token string) error {
	ctx, cancel := DefaultTimeout()
	defer cancel()

	_, err := r.refreshTokensCollection.DeleteOne(ctx, bson.M{"token": token})
	return err
}

//revokes all tokens for a user
func (r *mongoTokenRepository) RevokeAllUserTokens(userID string) error {
	ctx, cancel := DefaultTimeout()
	defer cancel()

	_, err := r.accessTokensCollection.DeleteMany(ctx, bson.M{"user_id": userID})
	if err != nil {
		return err
	}

	_, err = r.refreshTokensCollection.DeleteMany(ctx, bson.M{"user_id": userID})
	return err
} 