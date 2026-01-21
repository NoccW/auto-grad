package services

import (
	"auto-grad-backend/internal/models"
	"errors"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService struct {
	jwtSecret []byte
}

type Claims struct {
	UserID   uint   `json:"userId"`
	OpenID   string `json:"openId"`
	Role     string `json:"role"`
	UserRole string `json:"userRole"`
	jwt.StandardClaims
}

func NewAuthService(jwtSecret string) *AuthService {
	return &AuthService{
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *AuthService) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	claims := &Claims{
		UserID:   uint(user.ID),
		OpenID:   user.OpenID,
		Role:     user.Role,
		UserRole: "",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	if user.UserRole != nil {
		claims.UserRole = *user.UserRole
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
