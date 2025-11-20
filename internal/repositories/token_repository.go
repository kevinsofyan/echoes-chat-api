package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kevinsofyan/echoes-chat-api/internal/models"
	"gorm.io/gorm"
)

type TokenRepository interface {
	Create(ctx context.Context, token *models.Token) error
	FindByToken(ctx context.Context, tokenString string) (*models.Token, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Token, error)
	Delete(ctx context.Context, tokenString string) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) Create(ctx context.Context, token *models.Token) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *tokenRepository) FindByToken(ctx context.Context, tokenString string) (*models.Token, error) {
	var token models.Token
	err := r.db.WithContext(ctx).
		Where("token = ? AND expires_at > ?", tokenString, time.Now()).
		First(&token).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("token not found or expired")
		}
		return nil, err
	}

	return &token, nil
}

func (r *tokenRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]models.Token, error) {
	var tokens []models.Token
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND expires_at > ?", userID, time.Now()).
		Find(&tokens).Error

	return tokens, err
}

func (r *tokenRepository) Delete(ctx context.Context, tokenString string) error {
	return r.db.WithContext(ctx).
		Where("token = ?", tokenString).
		Delete(&models.Token{}).Error
}

func (r *tokenRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&models.Token{}).Error
}

func (r *tokenRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at <= ?", time.Now()).
		Delete(&models.Token{}).Error
}
