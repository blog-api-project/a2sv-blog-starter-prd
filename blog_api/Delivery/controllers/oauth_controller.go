package controllers

import (
	"blog_api/Delivery/dtos"
	contracts_usecases "blog_api/Domain/contracts/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)


type OAuthController struct {
	oauthUseCase contracts_usecases.IOAuthUseCase
}

func NewOAuthController(oauthUseCase contracts_usecases.IOAuthUseCase) *OAuthController {
	return &OAuthController{
		oauthUseCase: oauthUseCase,
	}
}

// starts the OAuth flow for a specific provider
func (oc *OAuthController) InitiateOAuthFlow(c *gin.Context) {
	provider := c.Param("provider")

	authURL, err := oc.oauthUseCase.InitiateOAuthFlow(provider)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"auth_url": authURL,
		"provider": provider,
	})
}

// HandleOAuthCallback processes the OAuth callback
func (oc *OAuthController) HandleOAuthCallback(c *gin.Context) {
	provider := c.Param("provider")
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code is required"})
		return
	}

	result, err := oc.oauthUseCase.HandleOAuthCallback(provider, code, state)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := dtos.OAuthLoginResponseDTO{
		Message:      "OAuth login successful",
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    900, // 15 minutes
		IsNewUser:    result.IsNewUser,
		User: dtos.UserResponseDTO{
			ID:        result.User.ID,
			Username:  result.User.Username,
			FirstName: result.User.FirstName,
			LastName:  result.User.LastName,
			Email:     result.User.Email,
		},
	}

	c.JSON(http.StatusOK, response)
}

// links an OAuth account to an existing user
func (oc *OAuthController) LinkOAuthToExistingUser(c *gin.Context) {
	provider := c.Param("provider")
	userID := c.GetString("user_id") 

	var linkDTO dtos.LinkOAuthDTO
	if err := c.ShouldBindJSON(&linkDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := oc.oauthUseCase.LinkOAuthToExistingUser(provider, linkDTO.Code, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "OAuth account linked successfully",
		"provider": provider,
	})
}
