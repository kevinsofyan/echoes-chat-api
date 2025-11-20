package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kevinsofyan/echoes-chat-api/internal/models"
	"github.com/kevinsofyan/echoes-chat-api/internal/repositories"
)

type MessageService interface {
	CreateMessage(ctx context.Context, req CreateMessageRequest) (*models.Message, error)
	GetMessageByID(ctx context.Context, id uuid.UUID) (*models.Message, error)
	GetMessagesByRoomID(ctx context.Context, roomID uuid.UUID, limit, offset int) ([]models.Message, error)
	UpdateMessage(ctx context.Context, id uuid.UUID, req UpdateMessageRequest) (*models.Message, error)
	DeleteMessage(ctx context.Context, id uuid.UUID) error
}

type CreateMessageRequest struct {
	RoomID    uuid.UUID  `json:"room_id" validate:"required"`
	SenderID  uuid.UUID  `json:"sender_id" validate:"required"`
	Content   string     `json:"content" validate:"required"`
	Type      string     `json:"type" validate:"required,oneof=text image file video audio"`
	FileURL   string     `json:"file_url,omitempty"`
	ReplyToID *uuid.UUID `json:"reply_to_id,omitempty"`
}

type UpdateMessageRequest struct {
	Content string `json:"content" validate:"required"`
}

type messageService struct {
	messageRepo repositories.MessageRepository
}

func NewMessageService(messageRepo repositories.MessageRepository) MessageService {
	return &messageService{
		messageRepo: messageRepo,
	}
}

func (s *messageService) CreateMessage(ctx context.Context, req CreateMessageRequest) (*models.Message, error) {
	message := &models.Message{
		RoomID:    req.RoomID,
		SenderID:  req.SenderID,
		Content:   req.Content,
		Type:      models.MessageType(req.Type),
		FileURL:   req.FileURL,
		ReplyToID: req.ReplyToID,
	}

	if err := s.messageRepo.Create(ctx, message); err != nil {
		return nil, err
	}

	// Fetch the message with preloaded relationships
	return s.messageRepo.FindByID(ctx, message.ID)
}

func (s *messageService) GetMessageByID(ctx context.Context, id uuid.UUID) (*models.Message, error) {
	return s.messageRepo.FindByID(ctx, id)
}

func (s *messageService) GetMessagesByRoomID(ctx context.Context, roomID uuid.UUID, limit, offset int) ([]models.Message, error) {
	if limit == 0 {
		limit = 50 // Default limit
	}
	return s.messageRepo.FindByRoomID(ctx, roomID, limit, offset)
}

func (s *messageService) UpdateMessage(ctx context.Context, id uuid.UUID, req UpdateMessageRequest) (*models.Message, error) {
	message, err := s.messageRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("message not found")
	}

	message.Content = req.Content
	message.IsEdited = true

	if err := s.messageRepo.Update(ctx, message); err != nil {
		return nil, err
	}

	return message, nil
}

func (s *messageService) DeleteMessage(ctx context.Context, id uuid.UUID) error {
	return s.messageRepo.Delete(ctx, id)
}
