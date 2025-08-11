package repositories

import "blog_api/Domain/models"

type IBlogRepository interface {
	CreateBlog(blog *models.Blog) error
	GetBlogs(query *models.BlogQuery) ([]models.Blog,int,error)
	UpdateBlog(blog models.Blog,BlogID string) ( *models.Blog,error)
	GetBlogByID(blogID string) (models.Blog,error)
	DeleteBlog( blogID string) error
	SearchBlogs(blogTitle string,authorID string)(*[]models.Blog,error)

	HasUserInteraction(userID, blogID, action string) (bool, error)
    AddUserInteraction(userID, blogID, action string) error
    RemoveUserInteraction(userID, blogID, action string) error
	IncrementLike(blogID string) error
    IncrementDislike(blogID string) error
	IncrementComment(blogID string)error
	DecrementComment(blogID string) error

}