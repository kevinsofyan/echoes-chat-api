package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kevinsofyan/echoes-chat-api/internal/models"
	"gorm.io/gorm"
)

type FriendshipRepository interface {
	Create(ctx context.Context, friendship *models.Friendship) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Friendship, error)
	FindByUsers(ctx context.Context, userID, friendID uuid.UUID) (*models.Friendship, error)
	Update(ctx context.Context, friendship *models.Friendship) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetFriendsList(ctx context.Context, userID uuid.UUID, status string, limit, offset int) ([]models.Friendship, error)
}

type friendshipRepository struct {
	db *gorm.DB
}

func NewFriendshipRepository(db *gorm.DB) FriendshipRepository {
	return &friendshipRepository{db: db}
}

func (r *friendshipRepository) Create(ctx context.Context, friendship *models.Friendship) error {
	return r.db.WithContext(ctx).Create(friendship).Error
}

func (r *friendshipRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Friendship, error) {
	var friendship models.Friendship
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Friend").
		First(&friendship, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("friendship not found")
		}
		return nil, err
	}

	return &friendship, nil
}

func (r *friendshipRepository) FindByUsers(ctx context.Context, userID, friendID uuid.UUID) (*models.Friendship, error) {
	var friendship models.Friendship
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Friend").
		Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
			userID, friendID, friendID, userID).
		First(&friendship).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not an error, just doesn't exist
		}
		return nil, err
	}

	return &friendship, nil
}

func (r *friendshipRepository) Update(ctx context.Context, friendship *models.Friendship) error {
	return r.db.WithContext(ctx).Save(friendship).Error
}

func (r *friendshipRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Friendship{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("friendship not found")
	}
	return nil
}

// GetFriendsList returns friendships filtered by status
func (r *friendshipRepository) GetFriendsList(ctx context.Context, userID uuid.UUID, status string, limit, offset int) ([]models.Friendship, error) {
	var friendships []models.Friendship
	query := r.db.WithContext(ctx).
		Preload("User").
		Preload("Friend")

	// Filter by user involvement
	query = query.Where("user_id = ? OR friend_id = ?", userID, userID)

	// Filter by status if specified
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&friendships).Error
	return friendships, err
}
