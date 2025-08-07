package repositories

import "blog_api/Domain/models"

type IBlogRepository interface {
	CreateBlog(blog *models.Blog) error
}