package usecases

import (
	"blog_api/Domain/contracts/repositories"
	"blog_api/Domain/models"
	"errors"
	"strings"
	"time"
)

type BlogUseCase struct {
	BlogRepo repositories.IBlogRepository
	UserRepo repositories.IUserRepository
}

func NewBlogUseCase(blogRepo repositories.IBlogRepository,userRepo repositories.IUserRepository) *BlogUseCase {
	return &BlogUseCase{
		BlogRepo: blogRepo,
	    UserRepo: userRepo,}
}

func (uc *BlogUseCase) CreateBlog(blog *models.Blog, AuthorID string) error {
	blog.AuthorID = AuthorID
	if blog.Title == "" || blog.Content == "" || blog.AuthorID == "" {
		return errors.New("Please include all required fields")
	}
	blog.AuthorID = AuthorID
	blog.CreatedAt = time.Now()
	blog.UpdatedAt = time.Now()
	blog.LikeCount = 0
	blog.DislikeCount = 0

	err := uc.BlogRepo.CreateBlog(blog)

	if err != nil {
		return err
	}
	return nil

}

func (uc *BlogUseCase) GetBlogs(query *models.BlogQuery)([]models.Blog, int,error){

	
	if query.Page <= 0{
		query.Page = 1
	}
	if query.PageSize <= 0{
		query.PageSize = 10
	}
	if query.SortBy == ""{
		query.SortBy = "recent"

	}
	blog,total,err := uc.BlogRepo.GetBlogs(query)
	if err != nil{
		return nil , 0,err
	}
	 
	return blog,total,err
}

func (uc *BlogUseCase) UpdateBlog(input *models.Blog, blogID string, authorID string) (*models.Blog, error) {
	blog, err := uc.BlogRepo.GetBlogByID(blogID)
	if err != nil {
		return nil, errors.New("unable to retrieve the blog")
	}

	if blog.AuthorID != authorID {
		return nil, errors.New("unauthorized access: you are not permitted to update this blog")
	}

	input.Title = strings.TrimSpace(input.Title)
	input.Content = strings.TrimSpace(input.Content)

	if input.Title == "" {
		return nil, errors.New("blog title must not be empty")
	}
	if input.Content == "" {
		return nil, errors.New("blog content must not be empty")
	}

	blog.Title = input.Title
	blog.Content = input.Content
	blog.Tags = input.Tags

	updatedBlog, err := uc.BlogRepo.UpdateBlog(blog, blogID)
	if err != nil {
		return nil, errors.New("failed to update the blog")
	}

	return updatedBlog, nil
}

func (uc *BlogUseCase) DeleteBlog(blogID string, authorID string) error {
	if strings.TrimSpace(blogID) == "" {
		return errors.New("invalid blog ID provided")
	}

	blog, err := uc.BlogRepo.GetBlogByID(blogID)
	if err != nil {
		return errors.New("blog not found")
	}

	if blog.AuthorID != authorID {
		return errors.New("unauthorized access: you are not permitted to delete this blog")
	}

	err = uc.BlogRepo.DeleteBlog(blogID)
	if err != nil {
		return errors.New("failed to delete the blog")
	}

	return nil
}

func (uc *BlogUseCase) SearchBlogs(searchQuery *models.BlogQuery)(*[]models.Blog,error) {
	var authorID string
	if searchQuery.Author != "" {
		user, err := uc.UserRepo.GetUserByUsername(searchQuery.Author)
		if err != nil {
			return nil, err
		}
		authorID = user.ID
	
	}

	blogs, err := uc.BlogRepo.SearchBlogs(searchQuery.Title,authorID)
	if err != nil {
		return nil, err
	}
	return blogs, nil
}

func (uc *BlogUseCase) LikeBlog(userID, blogID string) error {
   
    liked, err := uc.BlogRepo.HasUserInteraction(userID, blogID, "like")
    if err != nil {
        return err
    }
    if liked {
        return errors.New("user has already liked this blog")
    }
    err = uc.BlogRepo.AddUserInteraction(userID, blogID, "like")
    if err != nil {
        return err
    }

    err = uc.BlogRepo.IncrementLike(blogID)
    if err != nil {
        if rbErr := uc.BlogRepo.RemoveUserInteraction(userID, blogID, "like"); rbErr != nil {
        }
        return err
    }
    return nil
}

func (uc *BlogUseCase) DislikeBlog(userID, blogID string) error{
	dislike,err := uc.BlogRepo.HasUserInteraction(userID,blogID,"dislike")
	if err != nil{
		return err
	}
	if dislike{
		return errors.New("user has already disliked this blog")

	}
	err = uc.BlogRepo.AddUserInteraction(userID,blogID,"dislike")
	if err != nil{
		return err
	}
	err = uc.BlogRepo.IncrementDislike(blogID)
	if err != nil {
        if rbErr := uc.BlogRepo.RemoveUserInteraction(userID, blogID, "like"); rbErr != nil {
        }
        return err
    }
    return nil
	
}