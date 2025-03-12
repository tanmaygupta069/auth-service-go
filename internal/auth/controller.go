package auth

import (
	"context"
	"net/http"

	pb "github.com/tanmaygupta069/auth-service-go/generated"
	"github.com/tanmaygupta069/auth-service-go/pkg/auth"
)

type AuthController struct {
	pb.UnimplementedAuthServiceServer
	service AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		service: NewAuthService(),
	}
}

func (s *AuthController) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return &pb.LoginResponse{
			Response: &pb.Response{
				Code:    http.StatusBadRequest,
				Message: "email and password are required fields",
			},
			Token: "",
		}, nil
	}

	if !IsValidEmail(req.Email) {
		return &pb.LoginResponse{
			Response: &pb.Response{
				Code:    http.StatusBadRequest,
				Message: "enter a valid email",
			},
			Token: "",
		}, nil
	}

	if !s.service.UserRegistered(req.Email) {
		return &pb.LoginResponse{
			Response: &pb.Response{
				Code:    http.StatusBadRequest,
				Message: "invalid email or password",
			},
			Token: "",
		}, nil
	}

	hashedPassword, err := s.service.GetHashedPassword(req.Email)
	if err != nil {
		return &pb.LoginResponse{
			Response: &pb.Response{
				Code:    http.StatusInternalServerError,
				Message: "error in hashing password",
			},
			Token: "",
		}, nil
	}
	if !s.service.CheckPassword(req.Password, hashedPassword) {
		return &pb.LoginResponse{
			Response: &pb.Response{
				Code:    http.StatusUnauthorized,
				Message: "incorrect password",
			},
			Token: "",
		}, nil
	}
	token, err := auth.GenerateToken(req.Email)
	if err != nil {
		return &pb.LoginResponse{
			Response: &pb.Response{
				Code:    http.StatusInternalServerError,
				Message: "error in generating token",
			},
			Token: "",
		}, nil
	}

	return &pb.LoginResponse{
		Response: &pb.Response{
			Code:    http.StatusOK,
			Message: http.StatusText(http.StatusOK),
		},
		Token: token,
	}, nil
}

func (s *AuthController) Signup(ctx context.Context, req *pb.SignupRequest) (*pb.SignupResponse, error) {
	if req.Email == "" || req.Password == "" {
		return &pb.SignupResponse{
			Response: &pb.Response{
				Code:    http.StatusBadRequest,
				Message: "email and password are required fields",
			},
		}, nil
	}

	if !IsValidEmail(req.Email) {
		return &pb.SignupResponse{
			Response: &pb.Response{
				Code:    http.StatusBadRequest,
				Message: "enter a valid email",
			},
		}, nil
	}

	if !IsValidPassword(req.Password) {
		return &pb.SignupResponse{
			Response: &pb.Response{
				Code:    http.StatusBadRequest,
				Message: "password must be greater than 8 letters",
			},
		}, nil
	}

	if s.service.UserRegistered(req.Email) {
		return &pb.SignupResponse{
			Response: &pb.Response{
				Code:    http.StatusBadRequest,
				Message: "user already registered login",
			},
		}, nil
	}

	hashedPassword, err := s.service.HashPassword(req.Password)
	if err != nil {
		return &pb.SignupResponse{
			Response: &pb.Response{
				Code:    http.StatusInternalServerError,
				Message: "error in hashing password",
			},
		}, nil
	}
	if err := s.service.RegisterUser(req.Email, hashedPassword); err != nil {
		return &pb.SignupResponse{
			Response: &pb.Response{
				Code:    http.StatusInternalServerError,
				Message: "error in registering user",
			},
		}, nil
	}

	return &pb.SignupResponse{
		Response: &pb.Response{
			Code:    http.StatusOK,
			Message: "user registered",
		},
	}, nil
}
