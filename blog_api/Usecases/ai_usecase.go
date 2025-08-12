package usecases

import (
	"context"
	"blog_api/Domain/contracts/services"
)

type AIUseCase struct {
	aiService services.AIService
}

func NewAIUseCase(aiService services.AIService) *AIUseCase {
	return &AIUseCase{aiService: aiService}
}

func (uc *AIUseCase) GenerateBlogPost(ctx context.Context, topic string) (string, error) {
	return uc.aiService.GenerateBlogPost(ctx, topic)
}

func (uc *AIUseCase) SuggestImprovements(ctx context.Context, content string) (string, error) {
	return uc.aiService.SuggestImprovements(ctx, content)
}
