package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tanmaygupta069/auth-service-go/config"
	"golang.org/x/crypto/bcrypt"
)

var cfg, _ = config.GetConfig()

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(cfg.ServerConfig.JwtSecret)

type AuthService interface {
	UserRegistered(email string) bool
	HashPassword(password string) (string, error)
	RegisterUser(email string, password string) error
	CheckPassword(password string, hashedPassword string) bool
	GetHashedPassword(email string) (string, error)
}

type AuthServiceImp struct {
	repo RepositoryInterface
}

func NewAuthService() AuthService {
	return &AuthServiceImp{
		repo: NewRepository(),
	}
}

func (r *AuthServiceImp) UserRegistered(email string) bool {
	return r.repo.IsUserRegistered(email)
}

func (r *AuthServiceImp) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (r *AuthServiceImp) RegisterUser(email string, password string) error {
	err := r.repo.SaveUser(email, password)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthServiceImp) CheckPassword(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (r *AuthServiceImp) GetHashedPassword(email string) (string, error) {
	hashedPassword, err := r.repo.GetHashedPassword(email)
	if err != nil {
		return "", err
	}
	return hashedPassword, nil
}

func GenerateToken(email string) (string, error) {
	expirationTime := time.Now().Add(time.Hour)
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
