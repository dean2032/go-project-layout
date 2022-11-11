package middlewares

import (
	"net/http"
	"strings"

	"github.com/dean2032/go-project-layout/services"
	"github.com/dean2032/go-project-layout/utils/logging"
	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware middleware for jwt authentication
type JWTAuthMiddleware struct {
	service *services.JWTAuthService
}

// NewJWTAuthMiddleware creates new jwt auth middleware
func NewJWTAuthMiddleware(
	service *services.JWTAuthService,
) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{
		service: service,
	}
}

// Setup sets up jwt auth middleware
func (m *JWTAuthMiddleware) Setup() {}

// Handler handles middleware functionality
func (m *JWTAuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		t := strings.Split(authHeader, " ")
		if len(t) == 2 {
			authToken := t[1]
			authorized, err := m.service.Authorize(authToken)
			if authorized {
				c.Next()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			logging.Error(err.Error())
			c.Abort()
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "you are not authorized",
		})
		c.Abort()
	}
}
