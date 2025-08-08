package models

import "time"

// Role constants for user roles
const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)


type Role struct {
	ID        string
	Role      string    //"admin", "user"
	CreatedAt time.Time
	UpdatedAt time.Time
}
