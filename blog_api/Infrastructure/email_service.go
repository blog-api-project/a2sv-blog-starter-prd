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
	resetURL := os.Getenv("RESET_URL") 

	if smtpHost == "" || smtpPort == "" || smtpUsername == "" || smtpPassword == "" || resetURL == "" {
		fmt.Printf("Password reset token for %s: %s\n", email, resetToken)
		return nil
	}

	port, err := strconv.Atoi(smtpPort)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %v", err)
	}

	// Email content with HTML formatting
	subject := "Password Reset Request"
	resetLink := fmt.Sprintf("%s?token=%s", resetURL, resetToken)

	body := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6;">
			<h2>Password Reset Request</h2>
			<p>Hello,</p>
			<p>You have requested to reset your password. Click the button below to reset it:</p>
			<p>
				<a href="%s" style="background-color: #4CAF50; color: white; padding: 10px 20px;
				text-decoration: none; border-radius: 5px;">Reset Password</a>
			</p>
			<p>If the button doesn’t work, copy and paste this link into your browser:</p>
			<p><a href="%s">%s</a></p>
			<p>This link will expire in 1 hour.</p>
			<p>If you didn't request this password reset, please ignore this email.</p>
			<br>
			<p>Best regards,<br>Your Blog Team</p>
		</body>
		</html>
	`, resetLink, resetLink, resetLink)

	// Full email message (with MIME type for HTML)
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s",
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

// sends a confirmation email when password is changed
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

	subject := "Password Changed Successfully"

	body := `
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6;">
			<h2>Password Changed Successfully</h2>
			<p>Hello,</p>
			<p>Your password has been successfully changed.</p>
			<p>If you didn't make this change, please contact support immediately.</p>
			<br>
			<p>Best regards,<br>Your Blog Team</p>
		</body>
		</html>
	`

	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n%s",
		smtpUsername, email, subject, body)

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	err = smtp.SendMail(fmt.Sprintf("%s:%d", smtpHost, port), auth, smtpUsername, []string{email}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
