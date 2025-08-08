package usecases

import (
	repositories "blog_api/Domain/contracts/repositories"
	services "blog_api/Domain/contracts/services"
	usecases "blog_api/Domain/contracts/usecases"
	"blog_api/Domain/models"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

type UserUseCase struct {
	userRepo      repositories.IUserRepository
	passwordSvc   services.IPasswordService
	jwtSvc        services.IJWTService
	validationSvc services.IValidationService
	emailSvc     services.IEmailService
	tokenUseCase  usecases.ITokenUseCase
}

func NewUserUseCase(
	userRepo repositories.IUserRepository,
	passwordSvc services.IPasswordService,
	jwtSvc services.IJWTService,
	validationSvc services.IValidationService,
	emailSvc services.IEmailService,
	tokenUseCase usecases.ITokenUseCase,
) *UserUseCase {
	return &UserUseCase{
		userRepo:      userRepo,
		passwordSvc:   passwordSvc,
		jwtSvc:        jwtSvc,
		validationSvc: validationSvc,
		emailSvc:      emailSvc,
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
//initiates the password reset process
func (uc *UserUseCase) ForgotPassword(email string) error {
	if err := uc.validationSvc.ValidateEmail(email); err != nil {
		return err
	}

	// Check if user exists
	user, err := uc.userRepo.GetUserByEmail(email)
	if err != nil {
		// Don't reveal if email exists or not for security
		return nil
	}

	if !user.IsActive {
		return errors.New("account is deactivated")
	}

	resetToken, err := generateResetToken()
	if err != nil {
		return err
	}

	// Set reset token and expiry
	user.ResetPasswordToken = resetToken
	expiresAt := time.Now().Add(1 * time.Hour)
	user.ResetPasswordExpires = &expiresAt
	user.UpdatedAt = time.Now()

	if err := uc.userRepo.UpdateUser(user); err != nil {
		return err
	}

	return uc.emailSvc.SendPasswordResetEmail(user.Email, resetToken)
}

// resets the user's password using the reset token
 func (uc *UserUseCase) ResetPassword(token, newPassword string) error {
  // Validate new password strength
  if err := uc.validationSvc.ValidatePassword(newPassword); err != nil {
    return err
  }

  // Get user by reset token
  user, err := uc.userRepo.GetUserByResetToken(token)
  if err != nil {
    return err
  }

  // Check if reset token is expired
  if user.ResetPasswordExpires == nil || time.Now().After(*user.ResetPasswordExpires) {
    return errors.New("reset token has expired")
  }

  // Hash new password
  hashedPassword, err := uc.passwordSvc.HashPassword(newPassword)
  if err != nil {
    return err
  }

  // Update user password and clear reset token
  user.Password = hashedPassword
  user.ResetPasswordToken = ""
  user.ResetPasswordExpires = nil
  user.UpdatedAt = time.Now()

  // Update user in database
  if err := uc.userRepo.UpdateUser(user); err != nil {
    return err
  }

  // Verify the password was updated correctly by retrieving the user
  updatedUser, err := uc.userRepo.GetUserByEmail(user.Email)
  if err != nil {
    return errors.New("failed to verify password update")
  }

  // Verify the password hash matches
  if !uc.passwordSvc.CheckPasswordHash(newPassword, updatedUser.Password) {
    return errors.New("password update verification failed")
  }

  if err := uc.tokenUseCase.RevokeAllUserTokens(user.ID); err != nil {
    return err
  }

  return uc.emailSvc.SendPasswordChangedEmail(user.Email)
}

//generates a secure random token for password reset
func generateResetToken() (string, error) {
	bytes := make([]byte, 32) 
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}