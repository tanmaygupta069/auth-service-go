package auth

import (
	"fmt"

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
	ValidateToken(token string)(bool,error)
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

func (r *AuthServiceImp)ValidateToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return false,nil
	}

	return true, nil
}

