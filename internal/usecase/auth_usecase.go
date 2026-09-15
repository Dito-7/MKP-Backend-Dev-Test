package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"mkp-cinema-ticketing/internal/config"
	"mkp-cinema-ticketing/internal/entity"
	"mkp-cinema-ticketing/internal/repository"
	"mkp-cinema-ticketing/pkg/hash"
	"mkp-cinema-ticketing/pkg/jwt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserAlreadyExists  = errors.New("a user with this email already exists")
	ErrUserNotFound       = errors.New("user not found")
)

type AuthUsecase interface {
	Login(ctx context.Context, req entity.LoginRequest) (*entity.LoginResponse, error)
	Register(ctx context.Context, req entity.RegisterRequest) (*entity.UserResponse, error)
	GetProfile(ctx context.Context, userID string) (*entity.UserResponse, error)
}

type authUsecase struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthUsecase(userRepo repository.UserRepository, cfg *config.Config) AuthUsecase {
	return &authUsecase{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (u *authUsecase) Login(ctx context.Context, req entity.LoginRequest) (*entity.LoginResponse, error) {
	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if !hash.CheckPassword(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	token, err := jwt.GenerateToken(user.ID, user.Email, string(user.Role), u.cfg.JWTSecret, u.cfg.JWTExpiresHour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate auth token: %w", err)
	}

	return &entity.LoginResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: int64(u.cfg.JWTExpiresHour / time.Second),
		User:      user.ToResponse(),
	}, nil
}

func (u *authUsecase) Register(ctx context.Context, req entity.RegisterRequest) (*entity.UserResponse, error) {
	existing, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}

	pwHash, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	role := req.Role
	if role == "" {
		role = entity.RoleCustomer
	}

	newUser := &entity.User{
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: pwHash,
		PhoneNumber:  req.PhoneNumber,
		Role:         role,
		IsActive:     true,
	}

	if err := u.userRepo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	res := newUser.ToResponse()
	return &res, nil
}

func (u *authUsecase) GetProfile(ctx context.Context, userID string) (*entity.UserResponse, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	res := user.ToResponse()
	return &res, nil
}
