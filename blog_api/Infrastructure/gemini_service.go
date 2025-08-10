package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"google.golang.org/grpc/status"
)

type GeminiConfig struct {
	Model       string   // Required
	MaxTokens   *int     // Changed to pointer (nil = use API default)
	Temperature *float32 // Changed to pointer
}

type GeminiService struct {
	client *genai.Client
	config GeminiConfig
}

func NewGeminiService(cfg GeminiConfig) (*GeminiService, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("GEMINI_API_KEY not set in environment")
	}

	client, err := genai.NewClient(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &GeminiService{
		client: client,
		config: cfg,
	}, nil
}

func (s *GeminiService) GenerateBlogPost(ctx context.Context, topic string) (string, error) {
	log.Printf("Generating blog post about: %s", topic)
	defer func(start time.Time) {
		log.Printf("Generation completed in %v", time.Since(start))
	}(time.Now())

	model := s.client.GenerativeModel(s.config.Model)
	model.Temperature = s.config.Temperature // Now using pointer directly
	
	if s.config.MaxTokens != nil {
		maxTokens := int32(*s.config.MaxTokens)
		model.MaxOutputTokens = &maxTokens // Assign pointer to int32
	}

	prompt := fmt.Sprintf(
		"Write a 300-word professional blog post about: %s\n"+
			"Format: Markdown with headings (##), bullet points, and 1-2 code blocks\n"+
			"Tone: Technical but accessible\n"+
			"Audience: Software developers", topic)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		if grpcStatus, ok := status.FromError(err); ok {
			return "", fmt.Errorf("API error [%s]: %w", grpcStatus.Code(), err)
		}
		return "", fmt.Errorf("generation failed: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("empty response from API")
	}

	// Check all parts for text content
	var generatedText strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			generatedText.WriteString(string(text))
		}
	}

	if generatedText.Len() == 0 {
		return "", errors.New("no valid text content in response")
	}

	return generatedText.String(), nil
}

func (s *GeminiService) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}
