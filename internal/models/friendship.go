package models

import (
	"time"

	"github.com/google/uuid"
)

type FriendshipStatus string

const (
	FriendshipPending  FriendshipStatus = "pending"
	FriendshipAccepted FriendshipStatus = "accepted"
	FriendshipRejected FriendshipStatus = "rejected"
	FriendshipBlocked  FriendshipStatus = "blocked"
)

type Friendship struct {
	ID           uuid.UUID        `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID       uuid.UUID        `gorm:"type:uuid;not null;index" json:"user_id"`
	FriendID     uuid.UUID        `gorm:"type:uuid;not null;index" json:"friend_id"`
	Status       FriendshipStatus `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`
	ActionUserID uuid.UUID        `gorm:"type:uuid;not null" json:"action_user_id"` // Who did the action
	CreatedAt    time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time        `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User       User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Friend     User `gorm:"foreignKey:FriendID;constraint:OnDelete:CASCADE" json:"friend,omitempty"`
	ActionUser User `gorm:"foreignKey:ActionUserID;constraint:OnDelete:CASCADE" json:"-"`
}

func (Friendship) TableName() string {
	return "friendships"
}
