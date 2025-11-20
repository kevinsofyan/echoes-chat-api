package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/kevinsofyan/echoes-chat-api/internal/models"
	"github.com/kevinsofyan/echoes-chat-api/internal/repositories"
)

type UserService interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetAllUsers(ctx context.Context, limit, offset int) ([]models.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (*models.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	SetOnlineStatus(ctx context.Context, id uuid.UUID, isOnline bool) error
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *userService) GetAllUsers(ctx context.Context, limit, offset int) ([]models.User, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.userRepo.GetAll(ctx, limit, offset)
}

func (s *userService) UpdateUser(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *userService) SetOnlineStatus(ctx context.Context, id uuid.UUID, isOnline bool) error {
	return s.userRepo.UpdateOnlineStatus(ctx, id, isOnline)
}
