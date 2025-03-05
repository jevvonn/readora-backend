package dto

type CreateGenreRequest struct {
	Name string `json:"name" validate:"required"`
}
