package usecases

import "context"

type IAIUseCase interface {
	GenerateBlogPost(ctx context.Context, topic string) (string, error)
	SuggestImprovements(ctx context.Context, content string) (string, error)
}
