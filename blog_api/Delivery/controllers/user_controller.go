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
	
	accessToken, refreshToken, err := uc.tokenUsecase.GenerateAndStoreTokens(user.ID, user.RoleID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}
	
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"access_token": accessToken,
		"refresh_token": refreshToken,
		"token_type": "Bearer",
		"expires_in": 3600, // 1 hour
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

 