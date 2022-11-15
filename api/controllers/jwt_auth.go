package controllers

import (
	"github.com/dean2032/go-project-layout/services"
	"github.com/dean2032/go-project-layout/utils/logging"
	"github.com/gin-gonic/gin"
)

// JWTAuthController struct
type JWTAuthController struct {
	service     *services.JWTAuthService
	userService *services.UserService
}

// NewJWTAuthController creates new controller
func NewJWTAuthController(
	service *services.JWTAuthService,
	userService *services.UserService,
) *JWTAuthController {
	return &JWTAuthController{
		service:     service,
		userService: userService,
	}
}

// SignIn signs in user
func (jwt *JWTAuthController) SignIn(c *gin.Context) {
	logging.Info("SignIn route called")
	// Currently not checking for username and password
	// Can add the logic later if necessary.
	user, _ := jwt.userService.GetOneUser(uint(1))
	token := jwt.service.CreateToken(user)
	OnSuccess(c, token)
}

// Register registers user
func (jwt *JWTAuthController) Register(c *gin.Context) {
	logging.Info("Register route called")
	OnSuccess(c, "register route")
}
