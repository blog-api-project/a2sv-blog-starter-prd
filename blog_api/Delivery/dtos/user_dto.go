package dtos

// user registration request
type UserRegistrationDTO struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

//user login request
type UserLoginDTO struct {
	EmailOrUsername string `json:"email_or_username"`
	Password        string `json:"password"`
}

// login response
type UserLoginResponseDTO struct {
	Message       string `json:"message"`
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	TokenType     string `json:"token_type"`
	ExpiresIn     int64  `json:"expires_in"`
}

// refresh token request
type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token"`
}

// token validation request
type ValidateTokenDTO struct {
	AccessToken string `json:"access_token"`
}

// logout request
type LogoutDTO struct {
	AccessToken string `json:"access_token"`
} 