package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"type:varchar(255);not null" json:"name"`
	Username      string    `gorm:"type:varchar(255);uniqueIndex;not null;unique" json:"username"`
	Email         string    `gorm:"type:varchar(255);uniqueIndex;not null;unique" json:"email"`
	Password      string    `gorm:"type:varchar(255);not null;" json:"password"`
	Role          string    `gorm:"type:userRole;not null;default:'USER'" json:"role"`
	EmailVerified bool      `gorm:"type:boolean;not null" json:"email_verified"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
