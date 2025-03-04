package models

import "github.com/gofiber/fiber/v2"

type ResponseItf interface {
	Success(ctx *fiber.Ctx, message string) error
	SetData(data any) *JSONResponseModel
}

type NoErrInterface interface{}

type JSONResponseModel struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

func NewResponseModel() ResponseItf {
	return &JSONResponseModel{}
}

func (r *JSONResponseModel) SetData(data any) *JSONResponseModel {
	r.Data = data
	return r
}

func (r *JSONResponseModel) Success(ctx *fiber.Ctx, message string) error {
	r.Message = message
	return ctx.Status(fiber.StatusOK).JSON(r)
}
