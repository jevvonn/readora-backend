package dto

import "mime/multipart"

type CreateBookRequest struct {
	Title       string `form:"title" validate:"required"`
	Description string `form:"description"`
	Author      string `form:"author" validate:"required"`
	PublishDate string `form:"publish_date" validate:"required"`

	PDFFile *multipart.FileHeader
}
