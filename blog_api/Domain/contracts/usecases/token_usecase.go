package usecases

import (
	"blog_api/Domain/models"
)


type ITokenUseCase interface {
	StoreTokens(accessToken *models.AccessToken, refreshToken *models.RefreshToken) error
	GenerateAndStoreTokens(userID, roleID string) (*models.AccessToken, *models.RefreshToken, error)
	//generates a new access token with proper expiration
	RefreshAccessToken(userID, roleID string) (*models.AccessToken, error)
	ValidateAccessToken(token string) (bool, error)
	ValidateRefreshToken(token string) (bool, error)
	RevokeAccessToken(token string) error
	RevokeRefreshToken(token string) error
	RevokeAllUserTokens(userID string) error
} 