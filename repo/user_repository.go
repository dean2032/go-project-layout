package repo

import (
	"github.com/dean2032/go-project-layout/utils/logging"
	"gorm.io/gorm"
)

// UserRepository database structure
type UserRepository struct {
	*Database
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *Database) *UserRepository {
	return &UserRepository{
		Database: db,
	}
}

// WithTx enables repository with transaction
func (r *UserRepository) WithTx(txHandle *gorm.DB) *UserRepository {
	if txHandle == nil {
		logging.Error("Transaction Database not found in gin context. ")
		return r
	}
	r.Database.DB = txHandle
	return r
}
