package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"google.golang.org/grpc/status"
)

// GeminiConfig holds configuration for the Gemini service
type GeminiConfig struct {
	Model       string   // Required
	MaxTokens   *int     // Optional: nil for API default
	Temperature *float32 // Optional: nil for API default
}

// GeminiService implements AIService using Google's Gemini API
type GeminiService struct {
	client *genai.Client
	config GeminiConfig
}

// NewGeminiService creates a new Gemini service instance
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

// GenerateBlogPost generates blog content based on a topic
func (s *GeminiService) GenerateBlogPost(ctx context.Context, topic string) (string, error) {
	log.Printf("Generating blog post about: %s", topic)
	defer func(start time.Time) {
		log.Printf("Generation completed in %v", time.Since(start))
	}(time.Now())

	return s.generateContent(ctx, fmt.Sprintf(
		"Write a 300-word professional blog post about: %s\n"+
			"Format: Markdown with headings (##), bullet points, and 1-2 code blocks\n"+
			"Tone: Technical but accessible\n"+
			"Audience: Software developers", topic))
}

// SuggestImprovements suggests improvements for existing content
func (s *GeminiService) SuggestImprovements(ctx context.Context, content string) (string, error) {
	log.Printf("Suggesting improvements for content")
	defer func(start time.Time) {
		log.Printf("Suggestion completed in %v", time.Since(start))
	}(time.Now())

	return s.generateContent(ctx, fmt.Sprintf(
		"Suggest improvements for the following content:\n\n%s\n\n"+
			"Focus on clarity, technical accuracy, and engagement.\n"+
			"Format: Markdown bullet points\n"+
			"Tone: Constructive and professional", content))
}

// Common content generation logic
func (s *GeminiService) generateContent(ctx context.Context, prompt string) (string, error) {
	model := s.client.GenerativeModel(s.config.Model)
	if s.config.Temperature != nil {
		model.Temperature = s.config.Temperature
	}
	if s.config.MaxTokens != nil {
		maxTokens := int32(*s.config.MaxTokens)
		model.MaxOutputTokens = &maxTokens
	}

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

	var generatedText strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			generatedText.WriteString(string(text))
		}
	}

	if generatedText.Len() == 0 {
		return "", errors.New("no valid text content in response")
	}

	return strings.TrimSpace(generatedText.String()), nil
}

// clean up geminai resource
func (s *GeminiService) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}
