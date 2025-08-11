package services

import "blog_api/Domain/models"


type IOAuthService interface {
    GetAuthURL(state string) (string, error)
    ExchangeCodeForToken(code string) (*models.OAuthToken, error)
    GetUserInfo(accessToken string) (*models.OAuthUserInfo, error)
    GetProviderName() string
}