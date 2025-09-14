package gemini

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/generative-ai-go/genai"
)

type GeminiClientImpl struct {
	Model *genai.GenerativeModel
}

// NewGeminiClient cria um novo cliente de IA, pronto para ser usado.
func NewGeminiClient(model *genai.GenerativeModel) *GeminiClientImpl {
	return &GeminiClientImpl{Model: model}
}

// GenerateContent implements gemini.GeminiClient.
func (client *GeminiClientImpl) GenerateContent(
	ctx context.Context,
	instructions string,
	history []*genai.Content,
	prompt string,
) (string, []*genai.Content, error) {
	gemini := client.Model
	gemini.ResponseMIMEType = "text/plain"

	session := gemini.StartChat()
	session.History = history
	gemini.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(instructions)},
	}

	response, err := session.SendMessage(context.Background(), genai.Text(prompt))

	if err != nil {
		slog.Error("Erro ao enviar mensagem para o Gemini", "erro", err)
		return "", nil, fmt.Errorf("falha na chamada da API Gemini: %w", err)
	}

	if len(response.Candidates) > 0 && len(response.Candidates[0].Content.Parts) > 0 {
		if text, ok := response.Candidates[0].Content.Parts[0].(genai.Text); ok {
			return string(text), session.History, nil
		}
	}

	return "", session.History, fmt.Errorf("resposta do Gemini não contém conteúdo válido")
}
