package usecases

import (
	"blog_api/Domain/contracts/repositories"
	"blog_api/Domain/models"
	"time"
	"errors"
)

type BlogUseCase struct {
	BlogRepo repositories.IBlogRepository
}

func NewBlogUseCase(blogRepo repositories.IBlogRepository) *BlogUseCase {
	return &BlogUseCase{
		BlogRepo: blogRepo}
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