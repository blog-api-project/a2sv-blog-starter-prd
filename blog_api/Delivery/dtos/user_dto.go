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
type ForgotPasswordDTO struct {
	Email string `json:"email"`
}
type ResetPasswordDTO struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

type ForgotPasswordResponseDTO struct {
	Message string `json:"message"`
}
type ResetPasswordResponseDTO struct {
	Message string `json:"message"`
} 
type ProfileUpdateDTO struct {
	FirstName     string `json:"first_name" binding:"omitempty,max=50"`
	LastName      string `json:"last_name" binding:"omitempty,max=50"`
	Bio           string `json:"bio" binding:"omitempty,max=500"`
	ProfilePicture string `json:"profile_picture" binding:"omitempty,url"`
	ContactInfo   string `json:"contact_info" binding:"omitempty,max=100"`
}
type ProfileResponseDTO struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	Bio           string `json:"bio"`
	ProfilePicture string `json:"profile_picture"`
	ContactInfo   string `json:"contact_info"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
} 
