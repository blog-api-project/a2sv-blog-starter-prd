package usecases

import (
	"time"
)


type ITokenUseCase interface {
	StoreTokens(userID, accessToken, refreshToken string, accessExpiresAt, refreshExpiresAt time.Time) error
	GenerateAndStoreTokens(userID, roleID string) (string, string, error)
	//generates a new access token with proper expiration
	RefreshAccessToken(userID, roleID string) (string, error)
	ValidateAccessToken(token string) (bool, error)
	ValidateRefreshToken(token string) (bool, error)
	RevokeAccessToken(token string) error
	RevokeRefreshToken(token string) error
	RevokeAllUserTokens(userID string) error
} 