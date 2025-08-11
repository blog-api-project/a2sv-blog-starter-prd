package controllers

import (
	"blog_api/Delivery/dtos"
	"blog_api/Domain/contracts/usecases"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
	commentUseCase usecases.ICommentUseCase
}

func NewCommentController( commUsecase usecases.ICommentUseCase) *CommentController{
	return &CommentController{
		commentUseCase: commUsecase,
	}

}

func (ct *CommentController) CreateComment(c *gin.Context){
	userID := c.GetString("user_id")
	blogID := c.Param("id")
	var comment dtos.CommentDTO
	log.Println("the userID and the blogID %s and %s ",userID,blogID)

	if err := c.ShouldBindJSON(&comment); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid request"})
		return
	}

	err := ct.commentUseCase.CreateComment(blogID,userID,comment.Content)
	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}
	c.JSON(http.StatusOK,gin.H{"Message":"Comment created successfully"})

}

func (ct *CommentController) UpdateComment(c *gin.Context){
	commentID := c.Param("id")
	// userID := c.Param("userID")
	var comment dtos.CommentDTO
	if err := c.ShouldBindJSON(&comment); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid request"})
		return
	}
	err := ct.commentUseCase.UpdateComment(commentID,comment.Content)
	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}
	c.JSON(http.StatusOK,gin.H{"Message":"Comment updated successfully"})


}

func (ct *CommentController) DeleteComment(c *gin.Context){
	commentID := c.Param("id")
	err := ct.commentUseCase.DeleteComment(commentID)
	if err != nil{
		c.JSON(http.StatusBadGateway,gin.H{"Error":err.Error()})
		return
	}
	c.JSON(http.StatusOK,gin.H{"Message":"Comment Deleted!"})

}