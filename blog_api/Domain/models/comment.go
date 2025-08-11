package models

import "time"

type Comment struct {
	ID        string
	BlogID    string
	UserID    string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}