package usecases

import (
	"blog_api/Domain/contracts/repositories"
	"errors"
)

type CommentUseCase struct {
	commentRepo repositories.ICommentRepository
	blogRepo  repositories.IBlogRepository
}

func NewCommentUseCases( comRepo repositories.ICommentRepository,blogRepo repositories.IBlogRepository) *CommentUseCase{
	return &CommentUseCase{
		commentRepo: comRepo,
		blogRepo: blogRepo,
	}
}

func (uc *CommentUseCase) CreateComment(blogID,userID,content string) error{
	if blogID == "" || userID == "" {
        return errors.New("blogID and userID are required")
    }
    if content == "" {
        return errors.New("content cannot be empty")
    }


	err := uc.commentRepo.CreateComment(blogID,userID,content)
	if err != nil{
		return err
	}
	err = uc.blogRepo.IncrementComment(blogID)
	if err != nil{
		return err
	}
	return nil

}

func (uc *CommentUseCase) UpdateComment(commentID string, content string) error {
	// Check if the comment exists
	err := uc.commentRepo.CheckCommentExist(commentID)
	if err != nil {
		return err
	}

	// Attempt to update the comment
	err = uc.commentRepo.UpdateComment(commentID, content)
	if err != nil {
		return errors.New("Editing comment failed")
	}

	return nil
}

func (uc *CommentUseCase) DeleteComment(commentID string) error {
    if commentID == "" {
        return errors.New("commentID cannot be empty")
    }

    if err := uc.commentRepo.CheckCommentExist(commentID); err != nil {
        return err
    }

	comment,err := uc.commentRepo.GetCommentByID(commentID)
	if err != nil{
		return errors.New("comment not found")
	}

	blogID := comment.BlogID
	err = uc.commentRepo.DeleteComment(commentID)
	if err != nil{
		return errors.New("comment deletion failed")
	}
	err = uc.blogRepo.DecrementComment(blogID)
	if err != nil{
		return err
	}
	return nil
}
