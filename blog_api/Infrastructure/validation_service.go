package infrastructure

import (
	services "blog_api/Domain/contracts/services"
	"regexp"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

type ValidationServiceImpl struct{}

func NewValidationService() services.IValidationService {
	return &ValidationServiceImpl{}
}

// ValidateEmail validates email format
func (v *ValidationServiceImpl) ValidateEmail(email string) error {
	if email == "" {
		return &ValidationError{Field: "email", Message: "email is required"}
	}
	
	// Check email format using regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return &ValidationError{Field: "email", Message: "invalid email format"}
	}
	
	return nil
}

// ValidatePassword validates password strength
func (v *ValidationServiceImpl) ValidatePassword(password string) error {
	if password == "" {
		return &ValidationError{Field: "password", Message: "password is required"}
	}
	if len(password) < 8 {
		return &ValidationError{Field: "password", Message: "password must be at least 8 characters long"}
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return &ValidationError{Field: "password", Message: "password must contain at least one uppercase letter"}
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return &ValidationError{Field: "password", Message: "password must contain at least one lowercase letter"}
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return &ValidationError{Field: "password", Message: "password must contain at least one number"}
	}
	
	return nil
}

// ValidateUsername validates username format
func (v *ValidationServiceImpl) ValidateUsername(username string) error {
	if username == "" {
		return &ValidationError{Field: "username", Message: "username is required"}
	}
	if len(username) < 3 {
		return &ValidationError{Field: "username", Message: "username must be at least 3 characters long"}
	}
	if len(username) > 30 {
		return &ValidationError{Field: "username", Message: "username must be less than 30 characters"}
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(username) {
		return &ValidationError{Field: "username", Message: "username can only contain letters, numbers, and underscores"}
	}
	
	return nil
} 

