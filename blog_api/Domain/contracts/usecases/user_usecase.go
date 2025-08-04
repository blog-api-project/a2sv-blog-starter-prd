package usecases

import (
	"blog_api/Domain/models"
)


type IUserUseCase interface {
	RegisterUser(user *models.User) error
	LoginUser(emailOrUsername, password string) (*models.User, error)
	LogoutUser(userID string) error
} 