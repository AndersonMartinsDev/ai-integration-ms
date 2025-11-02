package gemini

import (
	"context"
	"fmt"
	"log/slog"

	speech "cloud.google.com/go/speech/apiv1"
	"github.com/google/generative-ai-go/genai"

	// Importar a SDK do Google Cloud Speech

	speechpb "cloud.google.com/go/speech/apiv1/speechpb"
)

type GeminiClientImpl struct {
	Model  *genai.GenerativeModel
	client *speech.Client
}

// NewGeminiClient cria um novo cliente de IA, pronto para ser usado.
func NewGeminiClient(model *genai.GenerativeModel) *GeminiClientImpl {

	// client, err := speech.NewClient(context.Background(), option.WithCredentialsFile("C:\\development\\judite-back-infra\\ai-integration-ms\\speech-key.json"))
	client, err := speech.NewClient(context.Background())
	if err != nil {
		slog.Error("falha ao criar o cliente Google Speech ", "erro", err)
		panic("Error para iniciar serviço speech-to-text")
	}

	return &GeminiClientImpl{Model: model, client: client}
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

func (gemi *GeminiClientImpl) TranscribeAudio(ctx context.Context, audioBytes []byte) (string, error) {
	// A API do Meta geralmente retorna audio no formato OGG (Opus codec).
	// O Google Cloud STT é otimizado para o codec LINEAR16, mas o OGG tambem pode ser processado.

	req := &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_OGG_OPUS, // Escolha o ENCODING correto
			SampleRateHertz: 16000,                               // O Meta padroniza em 16kHz
			LanguageCode:    "pt-BR",                             // Defina o idioma
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Content{Content: audioBytes},
		},
	}

	resp, err := gemi.client.Recognize(ctx, req)
	if err != nil {
		return "", fmt.Errorf("falha na chamada STT: %w", err)
	}

	if len(resp.Results) == 0 || len(resp.Results[0].Alternatives) == 0 {
		return "", fmt.Errorf("nenhum resultado de transcrição encontrado")
	}

	// Retorna a primeira e mais confiavel alternativa
	return resp.Results[0].Alternatives[0].Transcript, nil
}
