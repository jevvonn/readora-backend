package dto

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/jevvonn/readora-backend/internal/domain/entity"
)

type CreateBookRequest struct {
	Title       string `form:"title" validate:"required"`
	Description string `form:"description"`
	Author      string `form:"author" validate:"required"`
	PublishDate string `form:"publish_date" validate:"required"`
	Genres      string `form:"genres"`

	PDFFile *multipart.FileHeader
}

type GetBooksQuery struct {
	Search    string `query:"search" validate:"omitempty,max=255"`
	Genre     string `query:"genre" validate:"omitempty,max=255"`
	Limit     int    `query:"limit" validate:"omitempty,numeric,min=1,max=100"`
	Page      int    `query:"page" validate:"omitempty,numeric,min=1"`
	SortBy    string `query:"sort_by" validate:"omitempty,oneof=title author publish_date"`
	SortOrder string `query:"sort_order" validate:"omitempty,oneof=asc desc"`
	OwnerID   string `query:"owner_id" validate:"omitempty"`
}

type GetBooksResponse struct {
	ID          uuid.UUID `gorm:"primaryKey" json:"id,omitempty"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	PublishDate time.Time `json:"publish_date,omitempty"`

	CoverImageKey string `json:"cover_image_key,omitempty"`
	CoverImageURL string `json:"cover_image_url,omitempty"`

	FileKey string `json:"file_key,omitempty"`
	FileURL string `json:"file_url,omitempty"`

	IsPublic       bool   `json:"is_public,omitempty"`
	BookFileStatus string `json:"book_file_status,omitempty"`

	OwnerID uuid.UUID   `json:"owner_id,omitempty"`
	Owner   entity.User `json:"owner,omitempty"`

	Genres []entity.Genre `json:"genres,omitempty"`
}
