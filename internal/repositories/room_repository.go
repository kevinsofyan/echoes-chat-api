package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/kevinsofyan/echoes-chat-api/internal/models"
	"gorm.io/gorm"
)

type RoomRepository interface {
	Create(ctx context.Context, room *models.Room) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Room, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Room, error)
	Update(ctx context.Context, room *models.Room) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindDirectRoom(ctx context.Context, user1ID, user2ID uuid.UUID) (*models.Room, error)
	AddMember(ctx context.Context, member *models.RoomMember) error
}

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(ctx context.Context, room *models.Room) error {
	return r.db.WithContext(ctx).Create(room).Error
}

func (r *roomRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Room, error) {
	var room models.Room
	err := r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Members.User").
		First(&room, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Room, error) {
	var rooms []models.Room
	err := r.db.WithContext(ctx).
		Joins("JOIN room_members ON room_members.room_id = rooms.id").
		Where("room_members.user_id = ?", userID).
		Preload("Creator").
		Preload("Members.User").
		Find(&rooms).Error
	return rooms, err
}

func (r *roomRepository) Update(ctx context.Context, room *models.Room) error {
	return r.db.WithContext(ctx).Save(room).Error
}

func (r *roomRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Room{}, id).Error
}

func (r *roomRepository) FindDirectRoom(ctx context.Context, user1ID, user2ID uuid.UUID) (*models.Room, error) {
	var room models.Room

	// Find rooms where both users are members and type is direct
	err := r.db.WithContext(ctx).
		Joins("JOIN room_members rm1 ON rm1.room_id = rooms.id AND rm1.user_id = ?", user1ID).
		Joins("JOIN room_members rm2 ON rm2.room_id = rooms.id AND rm2.user_id = ?", user2ID).
		Where("rooms.type = ?", models.RoomTypeDirect).
		Preload("Members.User").
		First(&room).Error

	if err != nil {
		return nil, err
	}

	return &room, nil
}

func (r *roomRepository) AddMember(ctx context.Context, member *models.RoomMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}
