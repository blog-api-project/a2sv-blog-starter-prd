package usecases

import "blog_api/Domain/models"


type IOAuthUseCase interface {
	InitiateOAuthFlow(provider string) (string, error)
	HandleOAuthCallback(provider, code, state string) (*models.OAuthLoginResult, error)
	LinkOAuthToExistingUser(provider, code, userID string) error
}
