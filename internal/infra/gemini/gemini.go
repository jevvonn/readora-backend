package gemini

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"github.com/jevvonn/readora-backend/config"
	"google.golang.org/api/option"
)

func NewGeminiModel() *genai.GenerativeModel {
	ctx := context.Background()
	conf := config.Load()
	client, err := genai.NewClient(ctx, option.WithAPIKey(conf.GeminiAPIKey))

	if err != nil {
		panic(err)
	}

	model := client.GenerativeModel("gemini-2.0-flash")
	return model
}
