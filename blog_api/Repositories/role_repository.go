package repositories

import (
	"blog_api/Domain/contracts/repositories"
	"blog_api/Domain/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


type MongoRoleRepository struct {
	collection *mongo.Collection
}


func NewMongoRoleRepository(collection *mongo.Collection) repositories.IRoleRepository {
	return &MongoRoleRepository{
		collection: collection,
	}
}

// retrieves a role by ID
func (r *MongoRoleRepository) GetRoleByID(roleID string) (*models.Role, error) {
    objectID, err := primitive.ObjectIDFromHex(roleID)
    if err != nil {
        return nil, err
    }

    filter := bson.M{"_id": objectID}
    var doc bson.M
    if err := r.collection.FindOne(context.TODO(), filter).Decode(&doc); err != nil {
        return nil, err
    }
    var role models.Role
    role.ID = objectID.Hex()
    if v, ok := doc["role"].(string); ok {
        role.Role = v
    }
    if v, ok := doc["created_at"].(primitive.DateTime); ok {
        role.CreatedAt = v.Time()
    }
    if v, ok := doc["updated_at"].(primitive.DateTime); ok {
        role.UpdatedAt = v.Time()
    }
    return &role, nil
}

// retrieves role ID by role name
func (r *MongoRoleRepository) GetRoleIDByName(roleName string) (string, error) {
    filter := bson.M{"role": roleName}
    var doc bson.M
    if err := r.collection.FindOne(context.TODO(), filter).Decode(&doc); err != nil {
        return "", err
    }
    if id, ok := doc["_id"].(primitive.ObjectID); ok {
        return id.Hex(), nil
    }
    return "", nil
}
