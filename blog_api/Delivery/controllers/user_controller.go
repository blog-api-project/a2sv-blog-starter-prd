package controllers

import (
	"blog_api/Delivery/dtos"
	services "blog_api/Domain/contracts/services"
	usecases "blog_api/Domain/contracts/usecases"
	"blog_api/Domain/models"
	"net/http"

	"github.com/gin-gonic/gin"
)


type UserController struct {
	userUsecase  usecases.IUserUseCase
	tokenUsecase usecases.ITokenUseCase
	jwtService   services.IJWTService
}

func NewUserController(userUsecase usecases.IUserUseCase, tokenUsecase usecases.ITokenUseCase, jwtService services.IJWTService) *UserController {
	return &UserController{
		userUsecase:  userUsecase,
		tokenUsecase: tokenUsecase,
		jwtService:   jwtService,
	}
}

func (uc *UserController) ChangeToDomain(userDTO *dtos.UserRegistrationDTO) *models.User {
	var user models.User
	user.Username = userDTO.Username
	user.FirstName = userDTO.FirstName
	user.LastName = userDTO.LastName
	user.Email = userDTO.Email
	user.Password = userDTO.Password
	return &user
}

func (uc *UserController) Register(c *gin.Context) {
	var userDTO dtos.UserRegistrationDTO
	err := c.ShouldBindJSON(&userDTO)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	
	user := uc.ChangeToDomain(&userDTO)
	err = uc.userUsecase.RegisterUser(user)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.IndentedJSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func (uc *UserController) Login(c *gin.Context) {
	var userDTO dtos.UserLoginDTO
	err := c.ShouldBindJSON(&userDTO)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	
	user, err := uc.userUsecase.LoginUser(userDTO.EmailOrUsername, userDTO.Password)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	
	accessTokenModel, refreshTokenModel, err := uc.tokenUsecase.GenerateAndStoreTokens(user.ID, user.RoleID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}
	
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"access_token": accessTokenModel.Token,
		"refresh_token": refreshTokenModel.Token,
		"token_type": "Bearer",
		"expires_in": 900, // 15 minutes
	})
}

func (uc *UserController) Logout(c *gin.Context) {
	var logoutDTO dtos.LogoutDTO
	err := c.ShouldBindJSON(&logoutDTO)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	claims, err := uc.jwtService.ValidateJWT(logoutDTO.AccessToken)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	
	userID, ok := claims["user_id"].(string)
	if !ok {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims"})
		return
	}

	err = uc.userUsecase.LogoutUser(userID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}
	
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
//handles forgot password requests
func (uc *UserController) ForgotPassword(c *gin.Context) {
	var forgotPasswordDTO dtos.ForgotPasswordDTO
	err := c.ShouldBindJSON(&forgotPasswordDTO)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	
	err = uc.userUsecase.ForgotPassword(forgotPasswordDTO.Email)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Always return success to prevent email enumeration
	c.IndentedJSON(http.StatusOK, gin.H{"message": "If the email exists, a password reset link has been sent"})
}

//handles password reset requests
func (uc *UserController) ResetPassword(c *gin.Context) {
	var resetPasswordDTO dtos.ResetPasswordDTO
	err := c.ShouldBindJSON(&resetPasswordDTO)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	
	err = uc.userUsecase.ResetPassword(resetPasswordDTO.Token, resetPasswordDTO.NewPassword)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}
func (uc *UserController) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	var updateDTO dtos.ProfileUpdateDTO
	if err := c.ShouldBindJSON(&updateDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}
	update := models.UserProfileUpdate{
		FirstName:      updateDTO.FirstName,
		LastName:       updateDTO.LastName,
		Bio:            updateDTO.Bio,
		ProfilePicture: updateDTO.ProfilePicture,
		ContactInfo:    updateDTO.ContactInfo,
	}
	updatedUser, err := uc.userUsecase.UpdateUserProfile(userID, &update)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := dtos.ProfileResponseDTO{
		ID:             updatedUser.ID,
		Username:       updatedUser.Username,
		FirstName:      updatedUser.FirstName,
		LastName:       updatedUser.LastName,
		Email:          updatedUser.Email,
		Bio:            updatedUser.Bio,
		ProfilePicture: updatedUser.ProfilePicture,
		ContactInfo:    updatedUser.ContactInfo,
		CreatedAt:      updatedUser.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      updatedUser.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	c.JSON(http.StatusOK, resp)
}



