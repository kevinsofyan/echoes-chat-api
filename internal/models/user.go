package models

import "time"

type User struct {
	BaseModel
	Username string `gorm:"uniqueIndex;not null;size:50" json:"username"`
	Email    string `gorm:"uniqueIndex;not null;size:100" json:"email"`
	Password string `gorm:"not null" json:"-"` // "-" means don't include in JSON
	FullName string `gorm:"size:100" json:"full_name"`
	Avatar   string `gorm:"size:255" json:"avatar"`
	IsOnline bool   `gorm:"default:false" json:"is_online"`

	LastSeen time.Time `json:"last_seen"`

	// Relationships
	Messages     []Message    `gorm:"foreignKey:SenderID" json:"messages,omitempty"`
	RoomMembers  []RoomMember `gorm:"foreignKey:UserID" json:"room_members,omitempty"`
	CreatedRooms []Room       `gorm:"foreignKey:CreatedBy" json:"created_rooms,omitempty"`
}

func (User) TableName() string {
	return "users"
}
