package repositories

import (
	"blog_api/Domain/models"
)

type ITokenRepository interface {
	StoreAccessToken(accessToken *models.AccessToken) error
	StoreRefreshToken(refreshToken *models.RefreshToken) error
	// validates if an access token exists and is not expired
	ValidateAccessToken(token string) (bool, error)
	//validates if a refresh token exists and is not expired
	ValidateRefreshToken(token string) (bool, error)
	RevokeAccessToken(token string) error
	RevokeRefreshToken(token string) error
	RevokeAllUserTokens(userID string) error
} 