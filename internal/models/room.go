package models

type RoomType string

const (
	RoomTypeDirect RoomType = "direct"
	RoomTypeGroup  RoomType = "group"
)

type Room struct {
	BaseModel
	Name        string   `gorm:"size:100" json:"name"`
	Type        RoomType `gorm:"type:varchar(20);not null;default:'group'" json:"type"`
	Description string   `gorm:"type:text" json:"description"`
	Avatar      string   `gorm:"size:255" json:"avatar"`
	CreatedBy   uint     `gorm:"not null" json:"created_by"`

	// Relationships
	Creator  User         `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Members  []RoomMember `gorm:"foreignKey:RoomID" json:"members,omitempty"`
	Messages []Message    `gorm:"foreignKey:RoomID" json:"messages,omitempty"`
}

func (Room) TableName() string {
	return "rooms"
}
