package middlewares

import (
	"net/http"

	"github.com/dean2032/go-project-layout/constants"
	"github.com/dean2032/go-project-layout/repo"
	"github.com/dean2032/go-project-layout/utils/logging"
	"github.com/gin-gonic/gin"
)

// DatabaseTx middleware for transactions support for database
type DatabaseTx struct {
	handler *RequestHandler
	db      *repo.Database
}

// NewDatabaseTx creates new database transactions middleware
func NewDatabaseTx(
	handler *RequestHandler,
	db *repo.Database,
) *DatabaseTx {
	return &DatabaseTx{
		handler: handler,
		db:      db,
	}
}

// statusInList function checks if context writer status is in provided list
func statusInList(status int, statusList []int) bool {
	for _, i := range statusList {
		if i == status {
			return true
		}
	}
	return false
}

// Setup sets up database transaction middleware
func (m *DatabaseTx) Setup() {
	logging.Info("setting up database transaction middleware")

	m.handler.Gin.Use(func(c *gin.Context) {
		txHandle := m.db.DB.Begin()
		logging.Info("beginning database transaction")

		defer func() {
			if r := recover(); r != nil {
				txHandle.Rollback()
			}
		}()

		c.Set(constants.DBTransaction, txHandle)
		c.Next()

		// rollback transaction on server errors
		if c.Writer.Status() == http.StatusInternalServerError {
			logging.Info("rolling back transaction due to status code: 500")
			txHandle.Rollback()
		}

		// commit transaction on success status
		if statusInList(c.Writer.Status(), []int{http.StatusOK, http.StatusCreated}) {
			logging.Info("committing transactions")
			if err := txHandle.Commit().Error; err != nil {
				logging.Errorf("tx commit error: %s", err.Error())
			}
		}
	})
}
