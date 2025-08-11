package repositories

import "blog_api/Domain/models"

type IBlogRepository interface {
	CreateBlog(blog *models.Blog) error
	GetBlogs(query *models.BlogQuery) ([]models.Blog,int,error)
	UpdateBlog(blog models.Blog,BlogID string) ( *models.Blog,error)
	GetBlogByID(blogID string) (models.Blog,error)
	DeleteBlog( blogID string) error
}