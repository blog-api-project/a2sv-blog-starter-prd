package usecases


type ICommentUseCase interface {
	CreateComment(blogID, userID, content string) (error)
	UpdateComment(commentID,content string)(error)
	DeleteComment(commetnID string)(error)
}