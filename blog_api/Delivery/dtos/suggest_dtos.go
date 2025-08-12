package dtos

type SuggestImprovementsRequest struct {
	Content string `json:"content" binding:"required,min=50"`
}

type SuggestionsResponse struct {
	OriginalContent string `json:"original_content"`
	Suggestions     string `json:"suggestions"`
	GeneratedAt     string `json:"generated_at"`
}
