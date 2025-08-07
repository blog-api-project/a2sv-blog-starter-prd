package usecases

import "blog_api/Domain/models"

type IBlogUseCase interface {
	CreateBlog(blog *models.Blog, authorID string) (error)
}