package services

import (
	"errors"

	"github.com/dean2032/go-project-layout/config"
	"github.com/dean2032/go-project-layout/models"
	"github.com/dean2032/go-project-layout/utils/logging"
	"github.com/dgrijalva/jwt-go"
)

// JWTAuthService service relating to authorization
type JWTAuthService struct {
	cfg *config.Config
}

// NewJWTAuthService creates a new auth service
func NewJWTAuthService(cfg *config.Config) *JWTAuthService {
	return &JWTAuthService{
		cfg: cfg,
	}
}

// Authorize authorizes the generated token
func (s *JWTAuthService) Authorize(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTSecret), nil
	})
	if token.Valid {
		return true, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return false, errors.New("token malformed")
		}
		if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			return false, errors.New("token expired")
		}
	}
	return false, errors.New("couldn't handle token")
}

// CreateToken creates jwt auth token
func (s *JWTAuthService) CreateToken(user models.User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"name":  user.Name,
		"email": *user.Email,
	})

	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))

	if err != nil {
		logging.Errorf("JWT validation failed: %s", err.Error())
	}

	return tokenString
}
