package controllers

import (
	"net/http"
	"log"
	"os" 

	"github.com/gin-gonic/gin"
	"blog_api/Domain/contracts/services"
)

type AIController struct {
	aiService services.AIService
	logger    *log.Logger // Optional for structured logging
}

func NewAIController(aiService services.AIService) *AIController {
	return &AIController{
		aiService: aiService,
		logger:    log.New(os.Stdout, "[AI_CONTROLLER] ", log.LstdFlags),
	}
}

type GenerateRequest struct {
	Topic    string   `json:"topic" binding:"required,min=3,max=100"`
	Style    string   `json:"style,omitempty" enums:"professional,technical,casual"` // Optional
	Keywords []string `json:"keywords,omitempty" maxItems:"5"`                       // Optional
}

type GenerateResponse struct {
	Content    string   `json:"content"`
	Model      string   `json:"model,omitempty"`
	TokenUsage int      `json:"token_usage,omitempty"`
	Warnings   []string `json:"warnings,omitempty"`
}

func (c *AIController) GenerateBlogPost(ctx *gin.Context) {
	var req GenerateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.logger.Printf("Invalid request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":       "Invalid request format",
			"details":     err.Error(),
			"valid_types": []string{"professional", "technical", "casual"},
		})
		return
	}

	// Log the request (sanitize in production)
	c.logger.Printf("Generating content for topic: %s", req.Topic)

	content, err := c.aiService.GenerateBlogPost(ctx.Request.Context(), req.Topic)
	if err != nil {
		c.logger.Printf("Generation failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Content generation failed",
			"details": err.Error(), // Be cautious in production - don't expose internal errors
		})
		return
	}

	ctx.JSON(http.StatusOK, GenerateResponse{
		Content: content,
	})
}