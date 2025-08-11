package dtos


type OAuthLoginResponseDTO struct {
	Message      string           `json:"message"`
	AccessToken  string           `json:"access_token"`
	RefreshToken string           `json:"refresh_token"`
	TokenType    string           `json:"token_type"`
	ExpiresIn    int              `json:"expires_in"`
	IsNewUser    bool             `json:"is_new_user"`
	User         UserResponseDTO  `json:"user"`
}

// the request for linking OAuth account to existing user
type LinkOAuthDTO struct {
	Code string `json:"code" binding:"required"`
}

type UserResponseDTO struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
} 