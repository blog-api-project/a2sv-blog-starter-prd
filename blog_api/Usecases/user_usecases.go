package usecases

import (
	repositories "blog_api/Domain/contracts/repositories"
	services "blog_api/Domain/contracts/services"
	usecases "blog_api/Domain/contracts/usecases"
	"blog_api/Domain/models"
	"errors"
	"time"
)

type UserUseCase struct {
	userRepo      repositories.IUserRepository
	passwordSvc   services.IPasswordService
	jwtSvc        services.IJWTService
	validationSvc services.IValidationService
	tokenUseCase  usecases.ITokenUseCase
}

func NewUserUseCase(
	userRepo repositories.IUserRepository,
	passwordSvc services.IPasswordService,
	jwtSvc services.IJWTService,
	validationSvc services.IValidationService,
	tokenUseCase usecases.ITokenUseCase,
) *UserUseCase {
	return &UserUseCase{
		userRepo:      userRepo,
		passwordSvc:   passwordSvc,
		jwtSvc:        jwtSvc,
		validationSvc: validationSvc,
		tokenUseCase:  tokenUseCase,
	}
}

// registers a new user
func (uc *UserUseCase) RegisterUser(user *models.User) error {
	if err := uc.validationSvc.ValidateEmail(user.Email); err != nil {
		return err
	}
	if err := uc.validationSvc.ValidatePassword(user.Password); err != nil {
		return err
	}
	if err := uc.validationSvc.ValidateUsername(user.Username); err != nil {
		return err
	}
	exists, err := uc.userRepo.CheckEmailExists(user.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already exists")
	}
	exists, err = uc.userRepo.CheckUsernameExists(user.Username)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := uc.passwordSvc.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// Set default values
	user.IsActive = true
	user.EmailVerified = false
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return uc.userRepo.CreateUser(user)
}

// authenticates a user
func (uc *UserUseCase) LoginUser(emailOrUsername, password string) (*models.User, error) {
	var user *models.User
	var err error

	user, err = uc.userRepo.GetUserByEmail(emailOrUsername)
	if err != nil {
		user, err = uc.userRepo.GetUserByUsername(emailOrUsername)
		if err != nil {
			return nil, errors.New("invalid credentials")
		}
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	if !uc.passwordSvc.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// handles user logout business logic
func (uc *UserUseCase) LogoutUser(userID string) error {
	return uc.tokenUseCase.RevokeAllUserTokens(userID)
}
