package services

import "context"

type AIService interface {
    GenerateBlogPost(ctx context.Context, topic string) (string, error)
}