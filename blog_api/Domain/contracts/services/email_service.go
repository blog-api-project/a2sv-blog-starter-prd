package services

type IEmailService interface {
	SendPasswordResetEmail(email, resetToken string) error
	SendPasswordChangedEmail(email string) error
} 