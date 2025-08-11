package dtos

type CommentDTO struct {
    Content string `json:"content" binding:"required"`
}
