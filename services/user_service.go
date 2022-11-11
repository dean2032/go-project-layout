package services

import (
	"github.com/dean2032/go-project-layout/models"
	"github.com/dean2032/go-project-layout/repo"
	"gorm.io/gorm"
)

// UserService service layer
type UserService struct {
	repository *repo.UserRepository
}

// NewUserService creates a new userservice
func NewUserService(repository *repo.UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

// WithTx delegates transaction to repository database
func (s *UserService) WithTx(txHandle *gorm.DB) *UserService {
	s.repository = s.repository.WithTx(txHandle)
	return s
}

// GetOneUser gets one user
func (s *UserService) GetOneUser(id uint) (user models.User, err error) {
	return user, s.repository.Find(&user, id).Error
}

// GetAllUser get all the user
func (s *UserService) GetAllUser() (users []models.User, err error) {
	return users, s.repository.Find(&users).Error
}

// CreateUser call to create the user
func (s *UserService) CreateUser(user models.User) error {
	return s.repository.Create(&user).Error
}

// UpdateUser updates the user
func (s *UserService) UpdateUser(user models.User) error {
	return s.repository.Save(&user).Error
}

// DeleteUser deletes the user
func (s *UserService) DeleteUser(id uint) error {
	return s.repository.Delete(&models.User{}, id).Error
}
