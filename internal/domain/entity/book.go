package entity

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID          uuid.UUID `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"type:varchar(255);not null;index" json:"title"`
	Description string    `gorm:"type:varchar(255);not null;index" json:"description"`
	Author      string    `gorm:"type:varchar(255);not null;index" json:"author"`
	PublishDate time.Time `gorm:"type:timestamp;not null" json:"publish_date"`

	CoverImageKey string `gorm:"type:varchar(255);not null" json:"cover_image_key"`
	CoverImageURL string `gorm:"type:varchar(255);not null" json:"cover_image_url"`

	FileKey string `gorm:"type:varchar(255);not null" json:"file_key"`
	FileURL string `gorm:"type:varchar(255);not null" json:"file_url"`

	OwnerID uuid.UUID `gorm:"type:varchar(255);not null" json:"owner_id"`
	Owner   User      `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"owner"`

	Genres []Genre `gorm:"many2many:book_genres;" json:"genres"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
