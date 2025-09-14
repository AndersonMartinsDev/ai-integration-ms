package gemini

import (
	"context"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiConfig struct {
	Url     string
	Model   string
	Version string
	Key     string
}

var (
	GENAI_CONFIG GeminiConfig
	GEMINI_MODEL *genai.GenerativeModel
)

func (gem GeminiConfig) SetConfig() {
	client, err := genai.NewClient(context.Background(), option.WithAPIKey(gem.Key))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	model := client.GenerativeModel(gem.Model)
	model.SetTemperature(1)
	model.SetTopK(40)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(8192)

	GENAI_CONFIG = gem
	GEMINI_MODEL = model
}
