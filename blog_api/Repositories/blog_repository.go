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
	interactionCollection *mongo.Collection
}

func NewMongoBlogRepository(collection *mongo.Collection,interactionCol *mongo.Collection) repositories.IBlogRepository {
	return &MongoBlogRepository{
		blogCollection: collection,
	interactionCollection: interactionCol,}

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

func (r *MongoBlogRepository) SearchBlogs(title string, authorID string) (*[]models.Blog, error) {
    filter := bson.M{}

    if title != "" {
        filter["title"] = bson.M{"$regex": title, "$options": "i"}
    }

    if authorID != "" {
        objID, err := primitive.ObjectIDFromHex(authorID)
        if err != nil {
            return nil, err
        }
        filter["author_id"] = objID
    }

    cursor, err := r.blogCollection.Find(context.Background(), filter)
    if err != nil {
        return nil, err
    }

    var blogs []models.Blog
    if err := cursor.All(context.Background(), &blogs); err != nil {
        return nil, err
    }

    return &blogs, nil
}


func (bc *MongoBlogRepository) HasUserInteraction(userID,blogID,action string)(bool,error){
	userObjID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return false, err
    }

    blogObjID, err := primitive.ObjectIDFromHex(blogID)
    if err != nil {
        return false, err
    }
	filter := bson.M{
		"userid":userObjID,
		"blogid":blogObjID,
		"action":action, 

	}
	count,err := bc.interactionCollection.CountDocuments(context.TODO(),filter)
	if err != nil{
		return false,err
	}
	return count > 0, nil
}

func (bc *MongoBlogRepository) AddUserInteraction(userID, blogID, action string) error {
    userObjID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return err
    }

    blogObjID, err := primitive.ObjectIDFromHex(blogID)
    if err != nil {
        return err
    }

    interact := bson.M{
        "userid":    userObjID,
        "blogid":    blogObjID,
        "action":    action,
        "createdat": time.Now(),
    }
    _, err = bc.interactionCollection.InsertOne(context.TODO(), interact)
    return err
}

func (bc *MongoBlogRepository) RemoveUserInteraction(userID, blogID, action string) error {
    userObjID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return err
    }

    blogObjID, err := primitive.ObjectIDFromHex(blogID)
    if err != nil {
        return err
    }

    filter := bson.M{
        "userid": userObjID,
        "blogid": blogObjID,
        "action": action,
    }
    res, err := bc.interactionCollection.DeleteOne(context.TODO(), filter)
    if err != nil {
        return err
    }
    if res.DeletedCount == 0 {
        return errors.New("no interaction found")
    }
    return nil
}


func (m *MongoBlogRepository) IncrementLike(blogID string) error {
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$inc": bson.M{"likecount": 1}}
	_, err = m.blogCollection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (m *MongoBlogRepository) DecrementDislike(blogID string) error {
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID, "dislikecount": bson.M{"$gt": 0}}
	update := bson.M{"$inc": bson.M{"dislikecount": -1}}
	res, err := m.blogCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.New("dislike count already zero or blog not found")
	}
	return nil
}

func (bc *MongoBlogRepository) IncrementDislike(blogID string)error{
	objID,err := primitive.ObjectIDFromHex(blogID)
	if err != nil{
		return err
	}
	filter := bson.M{"_id":objID}
	update := bson.M{"$inc":bson.M{"dislikecount":1}}
	_,err = bc.blogCollection.UpdateOne(context.TODO(),filter,update)
	return nil
}

func (bc *MongoBlogRepository) IncrementComment(blogID string) error{
	objID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil{
		return nil
	}
	filter := bson.M{"_id":objID}
	update := bson.M{"$inc":bson.M{"commentcount":1}}
	_,err = bc.blogCollection.UpdateOne(context.TODO(),filter,update)
	return nil
}

func (bc *MongoBlogRepository) DecrementComment(blogID string) error {
    objID, err := primitive.ObjectIDFromHex(blogID)
    if err != nil {
        return err
    }

    filter := bson.M{"_id": objID}
    update := bson.M{"$inc": bson.M{"commentcount": -1}}

    _, err = bc.blogCollection.UpdateOne(context.TODO(), filter, update)
    return err
}
