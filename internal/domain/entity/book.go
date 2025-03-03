package entity

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID          uuid.UUID `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Description string    `gorm:"type:varchar(255);not null" json:"description"`
	Author      string    `gorm:"type:varchar(255);not null" json:"author"`
	PublishDate time.Time `gorm:"type:timestamp;not null" json:"publish_date"`

	CoverImageKey string `gorm:"type:varchar(255);not null" json:"cover_image_key"`
	CoverImageURL string `gorm:"type:varchar(255);not null" json:"cover_image_url"`

	FileKey string `gorm:"type:varchar(255);not null" json:"file_key"`
	FileURL string `gorm:"type:varchar(255);not null" json:"file_url"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
