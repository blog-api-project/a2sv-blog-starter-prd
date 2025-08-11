package dtos

type GenerateBlogRequest struct {
	Topic    string   `json:"topic" binding:"required,min=3,max=100"`
	Style    string   `json:"style,omitempty" enums:"professional,casual,academic"`
	Keywords []string `json:"keywords,omitempty" maxItems:"5"`
}

type BlogPostResponse struct {
	Content      string    `json:"content"`
	GeneratedAt  string    `json:"generated_at"` // RFC3339 format
	Model        string    `json:"model"`
	TokenUsage   int       `json:"token_usage"`
	TimeTakenMs  int       `json:"time_taken_ms"`
}