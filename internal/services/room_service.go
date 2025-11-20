package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/kevinsofyan/echoes-chat-api/internal/models"
	"github.com/kevinsofyan/echoes-chat-api/internal/repositories"
)

type RoomService interface {
	CreateRoom(ctx context.Context, req CreateRoomRequest) (*models.Room, error)
	GetRoomByID(ctx context.Context, id uuid.UUID) (*models.Room, error)
	GetUserRooms(ctx context.Context, userID uuid.UUID) ([]models.Room, error)
}

type CreateRoomRequest struct {
	Name        string    `json:"name" validate:"required"`
	Type        string    `json:"type" validate:"required,oneof=direct group"`
	Description string    `json:"description"`
	CreatedBy   uuid.UUID `json:"created_by"`
}

type roomService struct {
	roomRepo repositories.RoomRepository
}

func NewRoomService(roomRepo repositories.RoomRepository) RoomService {
	return &roomService{
		roomRepo: roomRepo,
	}
}

func (s *roomService) CreateRoom(ctx context.Context, req CreateRoomRequest) (*models.Room, error) {
	room := &models.Room{
		Name:        req.Name,
		Type:        models.RoomType(req.Type),
		Description: req.Description,
		CreatedBy:   req.CreatedBy,
	}

	if err := s.roomRepo.Create(ctx, room); err != nil {
		return nil, err
	}

	return s.roomRepo.FindByID(ctx, room.ID)
}

func (s *roomService) GetRoomByID(ctx context.Context, id uuid.UUID) (*models.Room, error) {
	return s.roomRepo.FindByID(ctx, id)
}

func (s *roomService) GetUserRooms(ctx context.Context, userID uuid.UUID) ([]models.Room, error) {
	return s.roomRepo.FindByUserID(ctx, userID)
}
