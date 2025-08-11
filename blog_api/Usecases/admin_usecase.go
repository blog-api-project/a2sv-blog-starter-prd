package usecases

import (
	"blog_api/Domain/contracts/repositories"
	"blog_api/Domain/contracts/usecases"
	"errors"
)


type AdminUseCase struct {
	userRepo repositories.IUserRepository
	roleRepo repositories.IRoleRepository
}


func NewAdminUseCase(userRepo repositories.IUserRepository, roleRepo repositories.IRoleRepository) usecases.IAdminUseCase {
	return &AdminUseCase{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// promotes a user to admin 
func (uc *AdminUseCase) PromoteUser(adminID, targetUserID string) error {
	if adminID == targetUserID {
		return errors.New("admin cannot promote themselves")
	}
	admin, err := uc.userRepo.GetUserByID(adminID)
	if err != nil {
		return errors.New("acting admin not found")
	}
    if admin.RoleID != "admin" {
        adminRole, err := uc.roleRepo.GetRoleByID(admin.RoleID)
        if err != nil || adminRole.Role != "admin" {
            return errors.New("only admins can promote users")
        }
    }

	user, err := uc.userRepo.GetUserByID(targetUserID)
	if err != nil {
		return errors.New("target user not found")
	}

    if user.RoleID == "admin" {
        return errors.New("user is already an admin")
    }
    if user.RoleID != "" {
        if userRole, err := uc.roleRepo.GetRoleByID(user.RoleID); err == nil && userRole.Role == "admin" {
            return errors.New("user is already an admin")
        }
    }

	if !user.IsActive {
		return errors.New("cannot promote inactive user")
	}
	adminRoleID, err := uc.roleRepo.GetRoleIDByName("admin")
	if err != nil {
		return errors.New("admin role not found")
	}

	return uc.userRepo.UpdateUserRole(targetUserID, adminRoleID)
}

// demotes an admin to user 
func (uc *AdminUseCase) DemoteUser(adminID, targetUserID string) error {
	if adminID == targetUserID {
		return errors.New("admin cannot demote themselves")
	}

	admin, err := uc.userRepo.GetUserByID(adminID)
	if err != nil {
		return errors.New("acting admin not found")
	}
    if admin.RoleID != "admin" {
        adminRole, err := uc.roleRepo.GetRoleByID(admin.RoleID)
        if err != nil || adminRole.Role != "admin" {
            return errors.New("only admins can demote users")
        }
    }

	user, err := uc.userRepo.GetUserByID(targetUserID)
	if err != nil {
		return errors.New("target user not found")
	}

    if user.RoleID == "admin" {
        // ok
    } else {
        userRole, err := uc.roleRepo.GetRoleByID(user.RoleID)
        if err != nil || userRole.Role != "admin" {
            return errors.New("user is not an admin")
        }
    }

	// Prevent demoting the last admin
	adminCount, err := uc.userRepo.GetAdminCount()
	if err != nil {
		return errors.New("failed to count admins")
	}
	if adminCount <= 1 {
		return errors.New("cannot demote the last admin")
	}

	// Get user role ID for demotion
	userRoleID, err := uc.roleRepo.GetRoleIDByName("user")
	if err != nil {
		return errors.New("user role not found")
	}
	
	return uc.userRepo.UpdateUserRole(targetUserID, userRoleID)
}
