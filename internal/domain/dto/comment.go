package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jevvonn/readora-backend/internal/domain/entity"
)

type CreateCommentRequest struct {
	Content string  `json:"content" validate:"required"`
	Rating  float64 `json:"rating" validate:"required,gte=1,lte=5,numeric"`
}

type GetCommentsQuery struct {
	Limit        int    `query:"limit" validate:"omitempty,numeric,min=1,max=100"`
	Page         int    `query:"page" validate:"omitempty,numeric,min=1"`
	SortBy       string `query:"sort_by" validate:"omitempty,oneof=created_at rating"`
	SortOrder    string `query:"sort_order" validate:"omitempty,oneof=asc desc"`
	BookId       string `query:"book_id" validate:"required"`
	TopCommentId string `query:"top_comment_id" validate:"omitempty"`
}

type GetCommentsResponse struct {
	ID        string      `json:"id"`
	Content   string      `json:"content"`
	Rating    float64     `json:"rating"`
	User      entity.User `json:"user"`
	BookId    uuid.UUID   `json:"book_id"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
