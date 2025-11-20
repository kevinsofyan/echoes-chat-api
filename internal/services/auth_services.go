package services

import (
	"context"
	"errors"
	"time"

	"github.com/kevinsofyan/echoes-chat-api/internal/models"
	"github.com/kevinsofyan/echoes-chat-api/internal/repositories"
	"github.com/kevinsofyan/echoes-chat-api/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req RegisterRequest) (*models.User, error)
	Login(ctx context.Context, req LoginRequest) (*models.User, string, error)
	Logout(ctx context.Context, token string) error
	ValidateToken(ctx context.Context, tokenString string) (*models.Token, error)
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserRequest struct {
	FullName string `json:"full_name"`
	Avatar   string `json:"avatar"`
}

type authService struct {
	userRepo  repositories.UserRepository
	tokenRepo repositories.TokenRepository
}

func NewAuthService(userRepo repositories.UserRepository, tokenRepo repositories.TokenRepository) AuthService {
	return &authService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

func (s *authService) Register(ctx context.Context, req RegisterRequest) (*models.User, error) {
	existingUser, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	existingUser, _ = s.userRepo.FindByUsername(ctx, req.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		FullName: req.FullName,
		IsOnline: false,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (s *authService) Login(ctx context.Context, req LoginRequest) (*models.User, string, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	tokenString, err := utils.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}

	token := &models.Token{
		UserID:    user.ID,
		Token:     tokenString,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.tokenRepo.Create(ctx, token); err != nil {
		return nil, "", errors.New("failed to store token")
	}

	s.userRepo.UpdateOnlineStatus(ctx, user.ID, true)

	user.Password = ""
	return user, tokenString, nil
}

func (s *authService) Logout(ctx context.Context, token string) error {
	tokenData, err := s.tokenRepo.FindByToken(ctx, token)
	if err != nil {
		return errors.New("invalid token")
	}

	if err := s.tokenRepo.DeleteByUserID(ctx, tokenData.UserID); err != nil {
		return errors.New("failed to logout")
	}

	s.userRepo.UpdateOnlineStatus(ctx, tokenData.UserID, false)

	return nil
}

func (s *authService) ValidateToken(ctx context.Context, tokenString string) (*models.Token, error) {
	token, err := s.tokenRepo.FindByToken(ctx, tokenString)
	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	return token, nil
}
