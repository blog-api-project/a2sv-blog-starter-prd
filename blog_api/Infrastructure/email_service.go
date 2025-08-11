package infrastructure

import (
	"fmt"
	"net/smtp"
	"os"
	"strconv"
)


type EmailService struct{}


func NewEmailService() *EmailService {
	return &EmailService{}
}

// sends a password reset email to the user
func (es *EmailService) SendPasswordResetEmail(email, resetToken string) error {
	// Get SMTP configuration from environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	if smtpHost == "" || smtpPort == "" || smtpUsername == "" || smtpPassword == "" {
		fmt.Printf("Password reset token for %s: %s\n", email, resetToken)
		return nil
	}
	port, err := strconv.Atoi(smtpPort)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %v", err)
	}

	// Create email content
	subject := "Password Reset Request"
	body := fmt.Sprintf(`
		Hello,
		
		You have requested to reset your password. Please click the link below to reset your password:
		
		http://localhost:3000/reset-password?token=%s
		
		This link will expire in 1 hour.
		
		If you didn't request this password reset, please ignore this email.
		
		Best regards,
		Your Blog Team
	`, resetToken)
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", 
		smtpUsername, email, subject, body)

	// Set up authentication
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	// Send email
	err = smtp.SendMail(fmt.Sprintf("%s:%d", smtpHost, port), auth, smtpUsername, []string{email}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

//sends a confirmation email when password is changed
func (es *EmailService) SendPasswordChangedEmail(email string) error {
	// Get SMTP configuration from environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	if smtpHost == "" || smtpPort == "" || smtpUsername == "" || smtpPassword == "" {
		fmt.Printf("Password changed confirmation sent to: %s\n", email)
		return nil
	}
	port, err := strconv.Atoi(smtpPort)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %v", err)
	}

	// Create email content
	subject := "Password Changed Successfully"
	body := `
		Hello,
		
		Your password has been successfully changed.
		
		If you didn't make this change, please contact support immediately.
		
		Best regards,
		Your Blog Team
	`

	// Create email message
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", 
		smtpUsername, email, subject, body)

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	// Send email
	err = smtp.SendMail(fmt.Sprintf("%s:%d", smtpHost, port), auth, smtpUsername, []string{email}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
} 