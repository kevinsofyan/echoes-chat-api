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
	DeleteRoom(ctx context.Context, id uuid.UUID) error
	GetOrCreateDirectRoom(ctx context.Context, user1ID, user2ID uuid.UUID) (*models.Room, error)
}

type CreateRoomRequest struct {
	Name        string      `json:"name" validate:"required"`
	Type        string      `json:"type" validate:"required,oneof=direct group"`
	Description string      `json:"description"`
	CreatedBy   uuid.UUID   `json:"created_by"`
	MemberIDs   []uuid.UUID `json:"member_ids"`
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

	// Add creator as owner
	if err := s.roomRepo.AddMember(ctx, &models.RoomMember{
		RoomID: room.ID,
		UserID: req.CreatedBy,
		Role:   models.RoleOwner,
	}); err != nil {
		return nil, err
	}

	// Add other members
	for _, memberID := range req.MemberIDs {
		if memberID != req.CreatedBy {
			if err := s.roomRepo.AddMember(ctx, &models.RoomMember{
				RoomID: room.ID,
				UserID: memberID,
				Role:   models.RoleMember,
			}); err != nil {
				return nil, err
			}
		}
	}

	return s.roomRepo.FindByID(ctx, room.ID)
}

func (s *roomService) GetRoomByID(ctx context.Context, id uuid.UUID) (*models.Room, error) {
	return s.roomRepo.FindByID(ctx, id)
}

func (s *roomService) GetUserRooms(ctx context.Context, userID uuid.UUID) ([]models.Room, error) {
	return s.roomRepo.FindByUserID(ctx, userID)
}

func (s *roomService) DeleteRoom(ctx context.Context, id uuid.UUID) error {
	return s.roomRepo.Delete(ctx, id)
}

func (s *roomService) GetOrCreateDirectRoom(ctx context.Context, user1ID, user2ID uuid.UUID) (*models.Room, error) {
	// Try to find existing direct room
	existingRoom, err := s.roomRepo.FindDirectRoom(ctx, user1ID, user2ID)
	if err == nil {
		return existingRoom, nil
	}

	// Create new direct room
	room := &models.Room{
		Name:        "", // Direct rooms don't need names
		Type:        models.RoomTypeDirect,
		Description: "",
		CreatedBy:   user1ID,
	}

	if err := s.roomRepo.Create(ctx, room); err != nil {
		return nil, err
	}

	// Add both users as members
	if err := s.roomRepo.AddMember(ctx, &models.RoomMember{
		RoomID: room.ID,
		UserID: user1ID,
		Role:   models.RoleMember,
	}); err != nil {
		return nil, err
	}

	if err := s.roomRepo.AddMember(ctx, &models.RoomMember{
		RoomID: room.ID,
		UserID: user2ID,
		Role:   models.RoleMember,
	}); err != nil {
		return nil, err
	}

	return s.roomRepo.FindByID(ctx, room.ID)
}
