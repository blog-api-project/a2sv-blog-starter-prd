package dtos

type GenerateBlogRequest struct {
	Topic string `json:"topic" binding:"required,min=3"`
}

type BlogPostResponse struct {
	Content      string `json:"content"`
	GeneratedAt  string `json:"generated_at"`
	Model        string `json:"model"`
	TimeTakenMs  int64  `json:"time_taken_ms,omitempty"`
}
