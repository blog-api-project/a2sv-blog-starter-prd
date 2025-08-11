package repositories

import (
	"blog_api/Domain/contracts/repositories"
	"blog_api/Domain/models"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommentRepository struct {
	commentCollection *mongo.Collection
}

func NewMongoCommentRepository (commCol *mongo.Collection) repositories.ICommentRepository{
	return &CommentRepository{
		commentCollection: commCol,

	}
}

func (r *CommentRepository) CreateComment(blogID,userID,content string)(error){
	objectID := primitive.NewObjectID()
	newComment := bson.M{
		"_id":objectID,
		"blogId":blogID,
		"userId":userID,
		"createdat":time.Now(),
		"updatedat":time.Now(),

	}
	_,err := r.commentCollection.InsertOne(context.Background(),newComment)

	if err != nil{
		return err
	}
	return nil

}

func (r *CommentRepository) GetCommentByID(commentID string) (models.Comment, error) {
    oid, err := primitive.ObjectIDFromHex(commentID)
    if err != nil {
        return models.Comment{}, fmt.Errorf("invalid commentID: %w", err)
    }

    filter := bson.M{"_id": oid}

    var model models.Comment
    err = r.commentCollection.FindOne(context.Background(), filter).Decode(&model)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return models.Comment{}, fmt.Errorf("comment not found")
        }
        return models.Comment{}, err
    }

    comment := models.Comment{
        ID:        model.ID,
        BlogID:    model.BlogID,
        UserID:    model.UserID,
        Content:   model.Content,
        CreatedAt: model.CreatedAt,
        UpdatedAt: model.UpdatedAt,
    }

    return comment, nil
}


func (r *CommentRepository) CheckCommentExist(commentID string) error {
	// Convert string to ObjectID
	oid, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return fmt.Errorf("invalid comment ID format: %w", err)
	}

	filter := bson.M{"_id": oid}

	// Perform the query
	err = r.commentCollection.FindOne(context.Background(), filter).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("comment with ID %s not found", commentID)
		}
		return err
	}

	return nil
}


func (r *CommentRepository) UpdateComment(commentID, content string) error {
    oid, err := primitive.ObjectIDFromHex(commentID)
    if err != nil {
        return err
    }

    update := bson.M{
        "$set": bson.M{
            "content":    content,
            "updated_at": time.Now(),
        },
    }

    _, err = r.commentCollection.UpdateOne(context.Background(), bson.M{"_id": oid}, update)
    if err != nil {
        return err
    }

    return nil
}

func (r *CommentRepository) DeleteComment(commentID string) error {
    oid, err := primitive.ObjectIDFromHex(commentID)
    if err != nil {
        return err 
    }

    result, err := r.commentCollection.DeleteOne(context.Background(), bson.M{"_id": oid})
    if err != nil {
        return err
    }
    if result.DeletedCount == 0 {
        return fmt.Errorf("comment not found")
    }
    return nil
}


