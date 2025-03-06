package entity

import (
	"time"

	"github.com/google/uuid"
)

type Reply struct {
	ID      uuid.UUID `gorm:"primaryKey" json:"id,omitempty"`
	Content string    `gorm:"type:text;not null" json:"content,omitempty"`

	ParentId uuid.UUID `gorm:"type:varchar(255);default:null" json:"parent,omitempty"`
	Replies  []Reply   `gorm:"foreignKey:ParentId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"reply,omitempty"`

	CommentId uuid.UUID `gorm:"type:varchar(255);not null" json:"comment_id,omitempty"`
	Comment   Comment   `gorm:"foreignKey:CommentId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"comment,omitempty"`

	UserId uuid.UUID `gorm:"type:varchar(255);not null" json:"user_id,omitempty"`
	User   User      `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
