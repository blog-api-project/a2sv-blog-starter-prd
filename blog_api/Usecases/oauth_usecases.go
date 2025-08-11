package usecases

import (
	"blog_api/Domain/contracts/repositories"
	"blog_api/Domain/contracts/services"
	"blog_api/Domain/contracts/usecases"
	"blog_api/Domain/models"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type OAuthUseCase struct {
	userRepo      repositories.IUserRepository
	oauthRepo     repositories.IOAuthRepository
	oauthServices map[string]services.IOAuthService
	tokenUseCase  usecases.ITokenUseCase
    roleRepo      repositories.IRoleRepository
}

func NewOAuthUseCase(
	userRepo repositories.IUserRepository,
	oauthRepo repositories.IOAuthRepository,
	oauthServices map[string]services.IOAuthService,
    tokenUseCase usecases.ITokenUseCase,
    roleRepo repositories.IRoleRepository,
) *OAuthUseCase {
	return &OAuthUseCase{
		userRepo:      userRepo,
		oauthRepo:     oauthRepo,
		oauthServices: oauthServices,
        tokenUseCase:  tokenUseCase,
        roleRepo:      roleRepo,
	}
}

// InitiateOAuthFlow starts the OAuth flow for a specific provider
func (uc *OAuthUseCase) InitiateOAuthFlow(provider string) (string, error) {
	oauthService, exists := uc.oauthServices[provider]
	if !exists {
		return "", fmt.Errorf("unsupported OAuth provider: %s", provider)
	}

	// For class project simplicity, no state
	authURL, err := oauthService.GetAuthURL("")
	if err != nil {
		return "", fmt.Errorf("failed to generate auth URL: %v", err)
	}

	return authURL, nil
}

// processes the OAuth callback and creates/links user
func (uc *OAuthUseCase) HandleOAuthCallback(provider, code, state string) (*models.OAuthLoginResult, error) {
	oauthService, exists := uc.oauthServices[provider]
	if !exists {
		return nil, fmt.Errorf("unsupported OAuth provider: %s", provider)
	}

	// Exchange code for token
	token, err := oauthService.ExchangeCodeForToken(code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %v", err)
	}

	userInfo, err := oauthService.GetUserInfo(token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}

	existingOAuthUser, err := uc.oauthRepo.GetOAuthUserByProviderID(provider, userInfo.ProviderID)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments)  {
		return nil, fmt.Errorf("failed to check existing OAuth user: %v", err)
	}

	var user *models.User
	var isNewUser bool

	if existingOAuthUser != nil {
		// get the linked user
		user, err = uc.userRepo.GetUserByID(existingOAuthUser.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get linked user: %v", err)
		}

		existingOAuthUser.AccessToken = token.AccessToken
		existingOAuthUser.RefreshToken = token.RefreshToken
		existingOAuthUser.UpdatedAt = time.Now()

		if err := uc.oauthRepo.UpdateOAuthUser(existingOAuthUser); err != nil {
			return nil, fmt.Errorf("failed to update OAuth user: %v", err)
		}
	} else {
		existingUser, err := uc.userRepo.GetUserByEmail(userInfo.Email)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments)  {
			return nil, fmt.Errorf("failed to check existing user: %v", err)
		}

		if existingUser != nil {
			user = existingUser
			isNewUser = false

		
			oauthUser := &models.OAuthUser{
				Provider:     provider,
				ProviderID:   userInfo.ProviderID,
				Email:        userInfo.Email,
				Name:         userInfo.Name,
				Picture:      userInfo.Picture,
				AccessToken:  token.AccessToken,
				RefreshToken: token.RefreshToken,
				UserID:       user.ID,
				ExpiresAt:    nil,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}

			if err := uc.oauthRepo.CreateOAuthUser(oauthUser); err != nil {
				return nil, fmt.Errorf("failed to create OAuth user: %v", err)
			}
		} else {
			
            user = &models.User{
                Email:     userInfo.Email,
                FirstName: userInfo.Name,
                Username:  userInfo.Email, 
                IsActive:  true,
                CreatedAt: time.Now(),
                UpdatedAt: time.Now(),
            }

            if uc.roleRepo != nil {
                if roleID, err := uc.roleRepo.GetRoleIDByName("user"); err == nil && roleID != "" {
                    user.RoleID = roleID
                } else {
                    user.RoleID = "user"
                }
            } else {
                user.RoleID = "user"
            }

			if err := uc.userRepo.CreateUser(user); err != nil {
				return nil, fmt.Errorf("failed to create user: %v", err)
			}

			oauthUser := &models.OAuthUser{
				Provider:     provider,
				ProviderID:   userInfo.ProviderID,
				Email:        userInfo.Email,
				Name:         userInfo.Name,
				Picture:      userInfo.Picture,
				AccessToken:  token.AccessToken,
				RefreshToken: token.RefreshToken,
				UserID:       user.ID,
				ExpiresAt:    nil,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}

			if err := uc.oauthRepo.CreateOAuthUser(oauthUser); err != nil {
				return nil, fmt.Errorf("failed to create OAuth user: %v", err)
			}

			isNewUser = true
		}
	}

	accessTokenModel, refreshTokenModel, err := uc.tokenUseCase.GenerateAndStoreTokens(user.ID, user.RoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %v", err)
	}

	return &models.OAuthLoginResult{
		User:         user,
		AccessToken:  accessTokenModel.Token,
		RefreshToken: refreshTokenModel.Token,
		IsNewUser:    isNewUser,
	}, nil
}

// links an OAuth account to an existing user
func (uc *OAuthUseCase) LinkOAuthToExistingUser(provider, code, userID string) error {
	oauthService, exists := uc.oauthServices[provider]
	if !exists {
		return fmt.Errorf("unsupported OAuth provider: %s", provider)
	}

	token, err := oauthService.ExchangeCodeForToken(code)
	if err != nil {
		return fmt.Errorf("failed to exchange code for token: %v", err)
	}

	userInfo, err := oauthService.GetUserInfo(token.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to get user info: %v", err)
	}

	existingOAuthUser, err := uc.oauthRepo.GetOAuthUserByProviderID(provider, userInfo.ProviderID)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return fmt.Errorf("failed to check existing OAuth user: %v", err)
	}

	if existingOAuthUser != nil {
		return fmt.Errorf("OAuth account already linked to another user")
	}

	
	oauthUser := &models.OAuthUser{
		Provider:     provider,
		ProviderID:   userInfo.ProviderID,
		Email:        userInfo.Email,
		Name:         userInfo.Name,
		Picture:      userInfo.Picture,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		UserID:       userID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return uc.oauthRepo.CreateOAuthUser(oauthUser)
}


