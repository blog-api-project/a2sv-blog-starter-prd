package usecases

import "blog_api/Domain/models"

type IBlogUseCase interface {
	CreateBlog(blog *models.Blog, authorID string) (error)
	GetBlogs(query *models.BlogQuery)([]models.Blog,int,error)
	UpdateBlog(updateblog *models.Blog,AuthorID string,BlogID string) (*models.Blog,error)
	DeleteBlog(BlogID string, AuthorID string) error
}