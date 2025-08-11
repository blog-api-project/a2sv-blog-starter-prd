 package repositories

import "blog_api/Domain/models"


type IOAuthRepository interface {
	CreateOAuthUser(oauthUser *models.OAuthUser) error
	GetOAuthUserByProviderID(provider, providerID string) (*models.OAuthUser, error)
	GetOAuthUserByEmail(provider, email string) (*models.OAuthUser, error)
	UpdateOAuthUser(oauthUser *models.OAuthUser) error
	LinkOAuthToUser(oauthUserID, userID string) error
}
