package models

type MessageType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypeImage MessageType = "image"
	MessageTypeFile  MessageType = "file"
	MessageTypeVideo MessageType = "video"
	MessageTypeAudio MessageType = "audio"
)

type Message struct {
	BaseModel
	RoomID    uint        `gorm:"not null;index" json:"room_id"`
	SenderID  uint        `gorm:"not null;index" json:"sender_id"`
	Content   string      `gorm:"type:text;not null" json:"content"`
	Type      MessageType `gorm:"type:varchar(20);not null;default:'text'" json:"type"`
	FileURL   string      `gorm:"size:255" json:"file_url,omitempty"`
	IsEdited  bool        `gorm:"default:false" json:"is_edited"`
	ReplyToID *uint       `gorm:"index" json:"reply_to_id,omitempty"`

	// Relationships
	Room    Room     `gorm:"foreignKey:RoomID" json:"room,omitempty"`
	Sender  User     `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
	ReplyTo *Message `gorm:"foreignKey:ReplyToID" json:"reply_to,omitempty"`
}

func (Message) TableName() string {
	return "messages"
}
