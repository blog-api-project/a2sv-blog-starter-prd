package controllers

import (
	"blog_api/Delivery/dtos"
	contracts_usecases "blog_api/Domain/contracts/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	adminUseCase contracts_usecases.IAdminUseCase
}

func NewAdminController(adminUseCase contracts_usecases.IAdminUseCase) *AdminController {
	return &AdminController{adminUseCase: adminUseCase}
}


func (ac *AdminController) PromoteUser(c *gin.Context) {
	adminID := c.GetString("user_id") 
	targetUserID := c.Param("userID")

	err := ac.adminUseCase.PromoteUser(adminID, targetUserID)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "only admins can promote users" {
			statusCode = http.StatusForbidden
		} else if err.Error() == "acting admin not found" || err.Error() == "target user not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dtos.PromoteUserResponseDTO{
		Message: "User promoted to admin successfully",
		UserID:  targetUserID,
		NewRole: "admin",
	})
}


func (ac *AdminController) DemoteUser(c *gin.Context) {
	adminID := c.GetString("user_id") 
	targetUserID := c.Param("userID")

	err := ac.adminUseCase.DemoteUser(adminID, targetUserID)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "only admins can demote users" {
			statusCode = http.StatusForbidden
		} else if err.Error() == "acting admin not found" || err.Error() == "target user not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dtos.DemoteUserResponseDTO{
		Message: "User demoted to user successfully",
		UserID:  targetUserID,
		NewRole: "user",
	})
}
