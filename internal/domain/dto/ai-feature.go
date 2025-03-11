package dto

import "github.com/google/generative-ai-go/genai"

type HighlightTextRequest struct {
	HighlightText string `json:"highlight_text" validate:"required"`
	Page          string `json:"page" validate:"required"`
}

type HighlightTextResponse struct {
	AIResponse genai.Text `json:"response"`
}
