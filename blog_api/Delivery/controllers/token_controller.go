package controllers

import (
	"blog_api/Delivery/dtos"
	services "blog_api/Domain/contracts/services"
	usecases "blog_api/Domain/contracts/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TokenController struct {
	tokenUseCase usecases.ITokenUseCase
	jwtService    services.IJWTService
}

func NewTokenController(
	tokenUseCase usecases.ITokenUseCase,
	jwtService services.IJWTService,
) *TokenController {
	return &TokenController{
		tokenUseCase: tokenUseCase,
		jwtService:    jwtService,
	}
}

// validates an access token
func (tc *TokenController) ValidateToken(c *gin.Context) {
	var validateDTO dtos.ValidateTokenDTO
	err := c.ShouldBindJSON(&validateDTO)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	valid, err := tc.tokenUseCase.ValidateAccessToken(validateDTO.AccessToken)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	if !valid {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Token not found or expired"})
		return
	}
	claims, err := tc.jwtService.ValidateJWT(validateDTO.AccessToken)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse token"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"valid": true,
		"claims": claims,
	})
}

// refreshes an access token using a refresh token
func (tc *TokenController) RefreshToken(c *gin.Context) {
	var refreshDTO dtos.RefreshTokenDTO
	err := c.ShouldBindJSON(&refreshDTO)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	valid, err := tc.tokenUseCase.ValidateRefreshToken(refreshDTO.RefreshToken)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	if !valid {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not found or expired"})
		return
	}

	claims, err := tc.jwtService.ValidateRefreshToken(refreshDTO.RefreshToken)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims"})
		return
	}

	role, _ := claims["role"].(string)
	newAccessToken, err := tc.tokenUseCase.RefreshAccessToken(userID, role)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
		"token_type":   "Bearer",
		"expires_in":   900, // 15 minutes
	})
} 