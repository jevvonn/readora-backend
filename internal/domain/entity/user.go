package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `gorm:"primaryKey" json:"id,omitempty"`
	Name          string    `gorm:"type:varchar(255);not null" json:"name,omitempty"`
	Username      string    `gorm:"type:varchar(255);index;not null;unique" json:"username,omitempty"`
	Email         string    `gorm:"type:varchar(255);index;not null;unique" json:"email,omitempty"`
	Password      string    `gorm:"type:varchar(255);not null;" json:"password,omitempty"`
	Role          string    `gorm:"type:userRole;not null;default:'USER'" json:"role,omitempty"`
	EmailVerified bool      `gorm:"type:boolean;not null" json:"email_verified,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
