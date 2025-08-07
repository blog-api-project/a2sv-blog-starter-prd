package repositories

import (
	repositories "blog_api/Domain/contracts/repositories"
	"blog_api/Domain/models"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoBlogRepository struct {
	blogCollection *mongo.Collection
}

func NewMongoBlogRepository(collection *mongo.Collection) repositories.IBlogRepository {
	return &MongoBlogRepository{blogCollection: collection}

}

func BuildFilterID (blogId string) (bson.M, error ){
	objID , err := primitive.ObjectIDFromHex(blogId)
	if err != nil{
		return nil,errors.New("Invalide id given")
	}
	return bson.M{"_id" :objID},nil
	

}


func (m *MongoBlogRepository) CreateBlog(blog *models.Blog) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID := primitive.NewObjectID()

	blogDoc := bson.M{
		"_id":           objectID,
		"authorid":     blog.AuthorID,
		"title":         blog.Title,
		"content":       blog.Content,
		"imageurl":     blog.ImageURL,
		"tags":          blog.Tags,
		"postedat":     blog.PostedAt,
		"likecount":    blog.LikeCount,
		"dislikecount": blog.DislikeCount,
		"commentcount": blog.CommentCount,
		"sharecount":   blog.ShareCount,
		"aisuggestion": blog.AISuggestion,
		"createdat":    blog.CreatedAt,
		"updatedat":    blog.UpdatedAt,
	}

	_, err := m.blogCollection.InsertOne(ctx, blogDoc)
	if err != nil {
		return err
	}

	blog.ID = objectID.Hex()
	return nil
}
func (m *MongoBlogRepository) GetBlogs(query *models.BlogQuery) ([]models.Blog, int, error) {
	skip := (query.Page - 1) * query.PageSize
	limit := query.PageSize	
	var sort bson.D
	switch  query.SortBy {
		case "popular":
			sort = bson.D{{Key: "LikeCount", Value: -1}}
		case "discussed":
			sort = bson.D{{Key: "CommentCount", Value: -1}}
		case "shared":
			sort = bson.D{{Key: "ShareCount", Value: -1}}
		case "oldest":
			sort = bson.D{{Key: "PostedAt", Value: 1}}
		default: 
			sort = bson.D{{Key: "PostedAt", Value: -1}}
		}

	findOptions := options.Find().
		SetSort(sort).
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	filter := bson.M{}
	cursor, err := m.blogCollection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, 0, err
	}

	var blogs []models.Blog
	if err := cursor.All(context.TODO(), &blogs); err != nil {
		return nil, 0, err
	}

	total, err := m.blogCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, 0, err
	}
	return blogs, int(total), nil
}

func (m *MongoBlogRepository) GetBlogByID(blogID string) (models.Blog, error) {
	var blog models.Blog

	objectID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return models.Blog{}, errors.New("invalid blog ID format")
	}

	err = m.blogCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&blog)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Blog{}, errors.New("blog not found")
		}
		return models.Blog{}, err
	}

	blog.ID = objectID.Hex()
	return blog, nil
}

func (m *MongoBlogRepository) UpdateBlog(updatedBlog models.Blog, blogID string) (*models.Blog, error) {
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, fmt.Errorf("invalid blog ID: %v", err)
	}

	update := bson.M{
		"$set": bson.M{
			"title":   updatedBlog.Title,
			"content": updatedBlog.Content,
			"tags":    updatedBlog.Tags,
			"updatedat":time.Now(),
		},
	}

	_ , err = m.blogCollection.UpdateByID(context.TODO(), objID, update)
	if err != nil {
		return nil, err
	}


	return &updatedBlog, nil
}

func (bc *MongoBlogRepository) DeleteBlog(BlogID string) error{
	filter,err := BuildFilterID(BlogID)

	if err != nil{
		return errors.New("Error while decoding the id")
	}
	_,err = bc.blogCollection.DeleteOne(context.TODO(),filter)
	if err != nil{
		return errors.New("Deletion Failed")

	}
	return nil
}