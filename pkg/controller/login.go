package controller

import (
	"fmt"
	"net/http"

	"shopping_api/pkg/dto"
	"shopping_api/pkg/service"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	service *service.LoginService
}

func NewLoginController(service *service.LoginService) *LoginController {
	return &LoginController{service: service}
}

func (c *LoginController) Login(ctx *gin.Context) {
	var loginReq dto.LoginRequest
	err := ctx.ShouldBindJSON(&loginReq)

	if err != nil {
		fmt.Print("Error parsing login request body:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not parse request body"})
		return
	}

	loggedInUser, err := c.service.Login(ctx, &loginReq)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, loggedInUser)
}

func (c *LoginController) SignUp(ctx *gin.Context) {
	var signupReq dto.SignupRequest
	err := ctx.ShouldBindJSON(&signupReq)

	if err != nil {
		fmt.Print("Error parsing login request body:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not parse request body"})
		return
	}

	err = c.service.SignUp(ctx, &signupReq)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
