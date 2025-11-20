package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/kevinsofyan/echoes-chat-api/internal/models"
	"gorm.io/gorm"
)

type MessageRepository interface {
	Create(ctx context.Context, message *models.Message) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Message, error)
	FindByRoomID(ctx context.Context, roomID uuid.UUID, limit, offset int) ([]models.Message, error)
	Update(ctx context.Context, message *models.Message) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(ctx context.Context, message *models.Message) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *messageRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Message, error) {
	var message models.Message
	err := r.db.WithContext(ctx).
		Preload("Sender").
		Preload("ReplyTo").
		First(&message, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *messageRepository) FindByRoomID(ctx context.Context, roomID uuid.UUID, limit, offset int) ([]models.Message, error) {
	var messages []models.Message
	query := r.db.WithContext(ctx).
		Where("room_id = ?", roomID).
		Preload("Sender").
		Preload("ReplyTo").
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&messages).Error
	return messages, err
}

func (r *messageRepository) Update(ctx context.Context, message *models.Message) error {
	return r.db.WithContext(ctx).Save(message).Error
}

func (r *messageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Message{}, id).Error
}
