package gateway

import (
	"context"

	"github.com/google/generative-ai-go/genai"
)

// GeminiClient define o contrato para a comunicação com a API do Gemini.
type GeminiClient interface {
	GenerateContent(
		ctx context.Context,
		instructions string,
		history []*genai.Content,
		prompt string,
	) (string, []*genai.Content, error)
}
