package dto

type CreateCommentRequest struct {
	Content string  `json:"content" validate:"required"`
	Rating  float64 `json:"rating" validate:"required,gte=1,lte=5,numeric"`
}
