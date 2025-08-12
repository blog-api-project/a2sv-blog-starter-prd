package controllers

import (
	"blog_api/Delivery/dtos"
	"blog_api/Domain/contracts/usecases"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AIController struct {
	aiUseCase usecases.IAIUseCase
	logger    *log.Logger
}

func NewAIController(aiUseCase usecases.IAIUseCase) *AIController {
	return &AIController{
		aiUseCase: aiUseCase,
		logger:    log.New(log.Writer(), "[AI_CONTROLLER] ", log.LstdFlags),
	}
}

func (ct *AIController) GenerateBlogPost(c *gin.Context) {
	var req dtos.GenerateBlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ct.logger.Printf("Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "Invalid request format",
			"details":     err.Error(),
			"valid_types": []string{"professional", "technical", "casual"},
		})
		return
	}

	ct.logger.Printf("Generating content for topic: %s", req.Topic)

	content, err := ct.aiUseCase.GenerateBlogPost(c.Request.Context(), req.Topic)
	if err != nil {
		ct.logger.Printf("Generation failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Content generation failed",
			"details": "Please try again later",
		})
		return
	}

	c.JSON(http.StatusOK, dtos.BlogPostResponse{
		Content:      content,
		GeneratedAt:  time.Now().Format(time.RFC3339),
		Model:        "gemini-1.5-flash",
		TimeTakenMs:  0,
	})
}

func (ct *AIController) SuggestImprovements(c *gin.Context) {
	var req dtos.SuggestImprovementsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ct.logger.Printf("Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	suggestions, err := ct.aiUseCase.SuggestImprovements(c.Request.Context(), req.Content)
	if err != nil {
		ct.logger.Printf("Suggestion failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not generate suggestions",
		})
		return
	}

	c.JSON(http.StatusOK, dtos.SuggestionsResponse{
		OriginalContent: req.Content,
		Suggestions:     suggestions,
		GeneratedAt:     time.Now().Format(time.RFC3339),
	})
}

func (ct *AIController) GenerateBlogContentForPost(c *gin.Context) {
	blogID := c.Param("id")
	userID := c.GetString("user_id")

	var req dtos.GenerateBlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ct.logger.Printf("Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	content, err := ct.aiUseCase.GenerateBlogPost(c.Request.Context(), req.Topic)
	if err != nil {
		ct.logger.Printf("Generation failed for blog %s: %v", blogID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not generate content",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"blog_id": blogID,
		"content": content,
		"user_id": userID,
	})
}
