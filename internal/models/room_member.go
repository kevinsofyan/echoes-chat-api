package models

import "time"

type RoomMemberRole string

const (
	RoleMember RoomMemberRole = "member"
	RoleAdmin  RoomMemberRole = "admin"
	RoleOwner  RoomMemberRole = "owner"
)

type RoomMember struct {
	BaseModel
	RoomID   uint           `gorm:"not null" json:"room_id"`
	UserID   uint           `gorm:"not null" json:"user_id"`
	Role     RoomMemberRole `gorm:"type:varchar(20);not null;default:'member'" json:"role"`
	JoinedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"joined_at"`

	// Relationships
	Room Room `gorm:"foreignKey:RoomID" json:"room,omitempty"`
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (RoomMember) TableName() string {
	return "room_members"
}

// Composite unique index
func (RoomMember) TableIndexes() []string {
	return []string{
		"idx_room_user:room_id,user_id,unique",
	}
}
