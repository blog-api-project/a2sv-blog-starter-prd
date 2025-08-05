package models

import (
	"time"
)


type Token struct {
	ID        string
	UserID    string
	Token     string
	Type      string // "access" or "refresh"
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AccessToken struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RefreshToken struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
} 