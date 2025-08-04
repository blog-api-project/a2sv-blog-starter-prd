package usecases

import (
	repositories "blog_api/Domain/contracts/repositories"
	services "blog_api/Domain/contracts/services"
	"time"
)

type TokenUseCase struct {
	tokenRepo repositories.ITokenRepository
	jwtSvc    services.IJWTService
}

func NewTokenUseCase(
	tokenRepo repositories.ITokenRepository,
	jwtSvc services.IJWTService,
) *TokenUseCase {
	return &TokenUseCase{
		tokenRepo: tokenRepo,
		jwtSvc:    jwtSvc,
	}
}

// stores access and refresh tokens for a user
func (uc *TokenUseCase) StoreTokens(userID, accessToken, refreshToken string, accessExpiresAt, refreshExpiresAt time.Time) error {
	err := uc.tokenRepo.StoreAccessToken(userID, accessToken, accessExpiresAt)
	if err != nil {
		return err
	}
	
	err = uc.tokenRepo.StoreRefreshToken(userID, refreshToken, refreshExpiresAt)
	if err != nil {
		return err
	}
	
	return nil
}

// generates and stores tokens with proper expiration times
func (uc *TokenUseCase) GenerateAndStoreTokens(userID, roleID string) (string, string, error) {
	accessToken, err := uc.jwtSvc.GenerateJWT(userID, roleID)
	if err != nil {
		return "", "", err
	}
	
	refreshToken, err := uc.jwtSvc.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}
	
	accessExpiresAt := time.Now().Add(15 * time.Minute) // 15 minutes
	refreshExpiresAt := time.Now().Add(time.Hour * 24 * 7) // 7 days
	
	err = uc.StoreTokens(userID, accessToken, refreshToken, accessExpiresAt, refreshExpiresAt)
	if err != nil {
		return "", "", err
	}
	
	return accessToken, refreshToken, nil
}

//generates a new access token with proper expiration
func (uc *TokenUseCase) RefreshAccessToken(userID, roleID string) (string, error) {
	accessToken, err := uc.jwtSvc.GenerateJWT(userID, roleID)
	if err != nil {
		return "", err
	}
	
	accessExpiresAt := time.Now().Add(15 * time.Minute)
	
	err = uc.tokenRepo.StoreAccessToken(userID, accessToken, accessExpiresAt)
	if err != nil {
		return "", err
	}
	
	return accessToken, nil
}

// validates an access token
func (uc *TokenUseCase) ValidateAccessToken(token string) (bool, error) {
	_, err := uc.jwtSvc.ValidateJWT(token)
	if err != nil {
		return false, err
	}
	return uc.tokenRepo.ValidateAccessToken(token)
}

//validates a refresh token
func (uc *TokenUseCase) ValidateRefreshToken(token string) (bool, error) {
	_, err := uc.jwtSvc.ValidateRefreshToken(token)
	if err != nil {
		return false, err
	}
	
	return uc.tokenRepo.ValidateRefreshToken(token)
}

// revokes a specific access token
func (uc *TokenUseCase) RevokeAccessToken(token string) error {
	return uc.tokenRepo.RevokeAccessToken(token)
}

//revokes a specific refresh token
func (uc *TokenUseCase) RevokeRefreshToken(token string) error {
	return uc.tokenRepo.RevokeRefreshToken(token)
}

// revokes all tokens for a user
func (uc *TokenUseCase) RevokeAllUserTokens(userID string) error {
	return uc.tokenRepo.RevokeAllUserTokens(userID)
} 