package entity

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID          uuid.UUID `gorm:"primaryKey" json:"id,omitempty"`
	Title       string    `gorm:"type:varchar(255);not null;index" json:"title,omitempty"`
	Description string    `gorm:"type:text;not null;index" json:"description,omitempty"`
	Author      string    `gorm:"type:varchar(255);not null;index" json:"author,omitempty"`
	PublishDate time.Time `gorm:"type:timestamp;not null" json:"publish_date,omitempty"`

	CoverImageURL string `gorm:"type:varchar(255);not null" json:"cover_image_url,omitempty"`

	FileKey  string `gorm:"type:varchar(255);not null" json:"file_key,omitempty"`
	FileURL  string `gorm:"type:varchar(255);not null" json:"file_url,omitempty"`
	FileType string `gorm:"type:varchar(255);not null" json:"file_type,omitempty"`

	IsPublic bool `gorm:"default:false;not null" json:"is_public,omitempty"`

	Rating float64 `gorm:"-" json:"rating"`

	OwnerID uuid.UUID `gorm:"type:varchar(255);not null" json:"owner_id,omitempty"`
	Owner   User      `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"owner,omitempty"`

	Genres           []Genre `gorm:"many2many:book_genres;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"genres,omitempty"`
	FileUploadStatus string  `gorm:"type:book_file_upload_status;not null;DEFAULT:'QUEUE';" json:"file_upload_status,omitempty"`
	FileAIStatus     string  `gorm:"type:book_file_ai_status;not null;DEFAULT:'QUEUE';" json:"file_ai_status,omitempty"`

	Comments []Comment `gorm:"foreignKey:BookId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"comments,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
