package repositories

import (
	repositories "blog_api/Domain/contracts/repositories"
	"blog_api/Domain/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoBlogRepository struct {
	blogCollection *mongo.Collection
}

func NewMongoBlogRepository(collection *mongo.Collection) repositories.IBlogRepository {
	return &MongoBlogRepository{blogCollection: collection}

}


func (m *MongoBlogRepository) CreateBlog(blog *models.Blog) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID := primitive.NewObjectID()

	blogDoc := bson.M{
		"_id":           objectID,
		"author_id":     blog.AuthorID,
		"title":         blog.Title,
		"content":       blog.Content,
		"image_url":     blog.ImageURL,
		"tags":          blog.Tags,
		"posted_at":     blog.PostedAt,
		"like_count":    blog.LikeCount,
		"dislike_count": blog.DislikeCount,
		"comment_count": blog.CommentCount,
		"share_count":   blog.ShareCount,
		"ai_suggestion": blog.AISuggestion,
		"created_at":    blog.CreatedAt,
		"updated_at":    blog.UpdatedAt,
	}

	_, err := m.blogCollection.InsertOne(ctx, blogDoc)
	if err != nil {
		return err
	}

	blog.ID = objectID.Hex()
	return nil
}