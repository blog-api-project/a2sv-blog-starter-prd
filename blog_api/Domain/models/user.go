package models

import (
	"time"
)

type User struct {
	ID                   string
	RoleID               string
	OAuthID              *string
	Username             string
	FirstName            string
	LastName             string
	Email                string
	Password             string
	Bio                  string
	ProfilePicture       string
	ContactInfo          string
	IsActive             bool
	EmailVerified        bool
	ResetPasswordToken   string
	ResetPasswordExpires *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
} 