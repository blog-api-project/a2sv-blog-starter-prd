package repositories

import "blog_api/Domain/models"

type ICommentRepository interface {
	CreateComment(blogID, userID, content string) error
	CheckCommentExist(CommentID string) error
	UpdateComment(CommentID, content string) error
	DeleteComment(commentID string) error
	GetCommentByID(commentID string) (models.Comment,error)
}