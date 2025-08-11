package repositories

import (
	"blog_api/Domain/models"
)

type IUserRepository interface {
	CreateUser(user *models.User) error
	GetUserByID(userID string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	CheckEmailExists(email string) (bool, error)
	CheckUsernameExists(username string) (bool, error)
	UpdateUser(user *models.User) error
	GetUserByResetToken(token string) (*models.User, error)
	UpdateUserRole(userID, newRole string) error
	GetAdminCount() (int, error)
	UpdateUserProfile(userID string, updateFields map[string]interface{}) error


} 