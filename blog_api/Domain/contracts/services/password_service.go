package services


type IPasswordService interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
} 