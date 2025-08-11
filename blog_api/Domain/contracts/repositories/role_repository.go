package repositories

import ("blog_api/Domain/models")


type IRoleRepository interface {
	GetRoleByID(roleID string) (*models.Role, error)
	GetRoleIDByName(roleName string) (string, error)
}
