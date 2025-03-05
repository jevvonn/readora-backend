package entity

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID      uuid.UUID `gorm:"primaryKey" json:"id,omitempty"`
	Content string    `gorm:"type:text;not null" json:"content,omitempty"`
	Rating  float64   `gorm:"type:float;not null" json:"rating,omitempty"`

	BookId uuid.UUID `gorm:"type:varchar(255);not null" json:"book_id,omitempty"`
	Book   Book      `gorm:"foreignKey:BookId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"book,omitempty"`

	UserId uuid.UUID `gorm:"type:varchar(255);not null" json:"user_id,omitempty"`
	User   User      `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
