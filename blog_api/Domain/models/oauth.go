package models

import "time"

// OAuth access token
type OAuthToken struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
	TokenType    string
}

//user information from OAuth provider
type OAuthUserInfo struct {
	ProviderID string
	Email      string
	Name       string
	Picture    string
}

// an OAuth user account linked to a platform user
type OAuthUser struct {
	ID           string
	UserID       string 
	Provider     string // "google", "github", "facebook"
	ProviderID   string // ID from the provider 
	Email        string
	Name         string
	Picture      string
	AccessToken  string
	RefreshToken string
	ExpiresAt    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// represents the result of an OAuth login
type OAuthLoginResult struct {
	User         *User
	AccessToken  string
	RefreshToken string
	IsNewUser    bool
}

// OAuthProvider constants
const (
	ProviderGoogle   = "google"
	ProviderGitHub   = "github"
	ProviderFacebook = "facebook"
)
