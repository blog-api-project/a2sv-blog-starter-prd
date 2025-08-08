package repositories

import (
	repositories "blog_api/Domain/contracts/repositories"
	"blog_api/Domain/models"
	"blog_api/Repositories/database"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


type mongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(collection *mongo.Collection) repositories.IUserRepository {
	return &mongoUserRepository{collection: collection}
}

//creates a new user in the database
func (r *mongoUserRepository) CreateUser(user *models.User) error {
	ctx, cancel := database.DefaultTimeout()
	defer cancel()

	// Generate ObjectID for the user
	objectID := primitive.NewObjectID()
	user.ID = objectID.Hex()

	// Convert domain model to BSON document
	doc := bson.M{
		"_id":                    objectID,
		"role_id":                user.RoleID,
		"oauth_id":               user.OAuthID,
		"username":               user.Username,
		"first_name":             user.FirstName,
		"last_name":              user.LastName,
		"email":                  user.Email,
		"password":               user.Password,
		"bio":                    user.Bio,
		"profile_picture":        user.ProfilePicture,
		"contact_info":           user.ContactInfo,
		"is_active":              user.IsActive,
		"email_verified":         user.EmailVerified,
		"reset_password_token":   user.ResetPasswordToken,
		"reset_password_expires": user.ResetPasswordExpires,
		"created_at":             user.CreatedAt,
		"updated_at":             user.UpdatedAt,
	}

	_, err := r.collection.InsertOne(ctx, doc)
	return err
}

//retrieves a user by ID
func (r *mongoUserRepository) GetUserByID(userID string) (*models.User, error) {
	ctx, cancel := database.DefaultTimeout()
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	var userData bson.M
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&userData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return r.documentToUser(userData)
}

// retrieves a user by email
func (r *mongoUserRepository) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := database.DefaultTimeout()
	defer cancel()

	var userData bson.M
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&userData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return r.documentToUser(userData)
}

//retrieves a user by username
func (r *mongoUserRepository) GetUserByUsername(username string) (*models.User, error) {
	ctx, cancel := database.DefaultTimeout()
	defer cancel()

	var userData bson.M
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&userData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return r.documentToUser(userData)
}

//checks if email already exists
func (r *mongoUserRepository) CheckEmailExists(email string) (bool, error) {
	ctx, cancel := database.DefaultTimeout()
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

//checks if username already exists
func (r *mongoUserRepository) CheckUsernameExists(username string) (bool, error) {
	ctx, cancel := database.DefaultTimeout()
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, bson.M{"username": username})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

//converts a BSON document to a User model
func (r *mongoUserRepository) documentToUser(userData bson.M) (*models.User, error) {
	user := &models.User{}

	if id, ok := userData["_id"].(primitive.ObjectID); ok {
		user.ID = id.Hex()
	}

	if roleID, ok := userData["role_id"].(string); ok {
		user.RoleID = roleID
	}

	if oauthID, ok := userData["oauth_id"].(*string); ok {
		user.OAuthID = oauthID
	}

	if username, ok := userData["username"].(string); ok {
		user.Username = username
	}

	if firstName, ok := userData["first_name"].(string); ok {
		user.FirstName = firstName
	}

	if lastName, ok := userData["last_name"].(string); ok {
		user.LastName = lastName
	}

	if email, ok := userData["email"].(string); ok {
		user.Email = email
	}

	if password, ok := userData["password"].(string); ok {
		user.Password = password
	}

	if bio, ok := userData["bio"].(string); ok {
		user.Bio = bio
	}

	if profilePicture, ok := userData["profile_picture"].(string); ok {
		user.ProfilePicture = profilePicture
	}

	if contactInfo, ok := userData["contact_info"].(string); ok {
		user.ContactInfo = contactInfo
	}

	if isActive, ok := userData["is_active"].(bool); ok {
		user.IsActive = isActive
	}

	if emailVerified, ok := userData["email_verified"].(bool); ok {
		user.EmailVerified = emailVerified
	}

	if resetPasswordToken, ok := userData["reset_password_token"].(string); ok {
		user.ResetPasswordToken = resetPasswordToken
	}

	if resetPasswordExpires, ok := userData["reset_password_expires"].(primitive.DateTime); ok {
		expiresTime := time.Unix(int64(resetPasswordExpires)/1000, 0)
		user.ResetPasswordExpires = &expiresTime
	}

	if createdAt, ok := userData["created_at"].(primitive.DateTime); ok {
		user.CreatedAt = time.Unix(int64(createdAt)/1000, 0)
	}

	if updatedAt, ok := userData["updated_at"].(primitive.DateTime); ok {
		user.UpdatedAt = time.Unix(int64(updatedAt)/1000, 0)
	}

	return user, nil
}

//updates an existing user in the database
func (r *mongoUserRepository) UpdateUser(user *models.User) error {
	ctx, cancel := database.DefaultTimeout()
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return errors.New("invalid user ID")
	}

	// Convert domain model to BSON document
	doc := bson.M{
		"role_id":                user.RoleID,
		"oauth_id":               user.OAuthID,
		"username":               user.Username,
		"first_name":             user.FirstName,
		"last_name":              user.LastName,
		"email":                  user.Email,
		"password":               user.Password,
		"bio":                    user.Bio,
		"profile_picture":        user.ProfilePicture,
		"contact_info":           user.ContactInfo,
		"is_active":              user.IsActive,
		"email_verified":         user.EmailVerified,
		"reset_password_token":   user.ResetPasswordToken,
		"reset_password_expires": user.ResetPasswordExpires,
		"updated_at":             user.UpdatedAt,
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": doc})
	return err
}

//retrieves a user by reset token
func (r *mongoUserRepository) GetUserByResetToken(token string) (*models.User, error) {
	ctx, cancel := database.DefaultTimeout()
	defer cancel()

	var userData bson.M
	err := r.collection.FindOne(ctx, bson.M{"reset_password_token": token}).Decode(&userData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("invalid or expired reset token")
		}
		return nil, err
	}

	return r.documentToUser(userData)
} 