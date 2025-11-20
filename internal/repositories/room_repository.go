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
		Preload("Members").
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
		Preload("Members").
		Find(&rooms).Error
	return rooms, err
}

func (r *roomRepository) Update(ctx context.Context, room *models.Room) error {
	return r.db.WithContext(ctx).Save(room).Error
}

func (r *roomRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Room{}, id).Error
}
