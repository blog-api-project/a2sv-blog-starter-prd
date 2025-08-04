package infrastructure

import (
	services "blog_api/Domain/contracts/services"

	"golang.org/x/crypto/bcrypt"
)

type PasswordServiceImpl struct{}

func NewPasswordService() services.IPasswordService {
	return &PasswordServiceImpl{}
}

//hashes a password using bcrypt
func (p *PasswordServiceImpl) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

//checks if a password matches a hash
func (p *PasswordServiceImpl) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
} 
