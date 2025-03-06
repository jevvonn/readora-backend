package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jevvonn/readora-backend/internal/domain/entity"
)

type CreateReplyRequest struct {
	ParentId string `json:"parent_id,omitempty"`
	Content  string `json:"content,omitempty" validate:"required"`
}

type GetRepliesQuery struct {
	Limit     int    `query:"limit" validate:"omitempty,numeric,min=1,max=100"`
	Page      int    `query:"page" validate:"omitempty,numeric,min=1"`
	SortOrder string `query:"sort_order" validate:"omitempty,oneof=asc desc"`
	ParentId  string `query:"parent_id" validate:"omitempty,oneof=asc desc"`
}

type GetRepliesResponse struct {
	ID        uuid.UUID   `json:"id"`
	Content   string      `json:"content"`
	ParentId  uuid.UUID   `json:"parent_id,omitempty"`
	CommentId uuid.UUID   `json:"comment_id"`
	User      entity.User `json:"user"`
	CreatedAt time.Time   `json:"created_at"`
}
