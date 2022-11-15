package controllers

import (
	"strconv"

	"github.com/dean2032/go-project-layout/constants"
	"github.com/dean2032/go-project-layout/models"
	"github.com/dean2032/go-project-layout/services"
	"github.com/dean2032/go-project-layout/utils/logging"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserController data type
type UserController struct {
	service *services.UserService
}

// NewUserController creates new user controller
func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		service: userService,
	}
}

// GetOneUser gets one user
func (u *UserController) GetOneUser(c *gin.Context) {
	paramID := c.Param("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		logging.Error(err.Error())
		OnError(c, err)
		return
	}
	user, err := u.service.GetOneUser(uint(id))

	if err != nil {
		logging.Error(err.Error())
		OnError(c, err)
		return
	}

	OnSuccess(c, user)
}

// GetUser gets the user
func (u *UserController) GetUser(c *gin.Context) {
	users, err := u.service.GetAllUser()
	if err != nil {
		logging.Error(err.Error())
	}
	OnSuccess(c, users)
}

// SaveUser saves the user
func (u *UserController) SaveUser(c *gin.Context) {
	user := models.User{}
	txHandle := c.MustGet(constants.DBTransaction).(*gorm.DB)

	if err := c.ShouldBindJSON(&user); err != nil {
		logging.Error(err.Error())
		OnError(c, err)
		return
	}

	if err := u.service.WithTx(txHandle).CreateUser(user); err != nil {
		logging.Error(err.Error())
		OnError(c, err)
		return
	}

	OnSuccess(c, nil)
}

// UpdateUser updates user
func (u *UserController) UpdateUser(c *gin.Context) {
	OnSuccess(c, nil)
}

// DeleteUser deletes user
func (u *UserController) DeleteUser(c *gin.Context) {
	paramID := c.Param("id")

	id, err := strconv.Atoi(paramID)
	if err != nil {
		logging.Error(err.Error())
		OnError(c, err)
		return
	}

	if err := u.service.DeleteUser(uint(id)); err != nil {
		logging.Error(err.Error())
		OnError(c, err)
		return
	}

	OnSuccess(c, nil)
}
