package models

import "github.com/gofiber/fiber/v2"

type ResponseItf interface {
	Success(ctx *fiber.Ctx, message string) error
	BadRequest(ctx *fiber.Ctx, err error, errData any) error
	InternalServerError(ctx *fiber.Ctx, err error, errData any) error
	SetData(data any) *JSONResponseModel
}

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

func (r *JSONResponseModel) BadRequest(ctx *fiber.Ctx, err error, errData any) error {
	res := &JSONResponseModel{
		Message: err.Error(),
		Errors:  errData,
	}
	return ctx.Status(fiber.StatusBadRequest).JSON(res)
}

func (r *JSONResponseModel) InternalServerError(ctx *fiber.Ctx, err error, errData any) error {
	res := &JSONResponseModel{
		Message: err.Error(),
		Errors:  errData,
	}
	return ctx.Status(fiber.StatusInternalServerError).JSON(res)
}
