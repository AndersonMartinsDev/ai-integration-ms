package services

import (
	"ai-integration-ms/internal/domain/model"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	// Assumindo que você tem acesso ao seu modelo e clientes
)

// Estrutura do MetaService (Com a configuracao do Meta)
type WWEBJSService struct {
	HttpClient    *http.Client
	WWEBJS_MS_URL string
}

func NewWWEBJSService(httpClient *http.Client) *WWEBJSService {
	return &WWEBJSService{
		HttpClient:    httpClient,
		WWEBJS_MS_URL: os.Getenv("WWEBJS_MS_URL"),
	}
}

func (c *WWEBJSService) DownloadAudioFile(ctx context.Context, modelRequest model.AiRequestModel) ([]byte, error) {
	// 1. Aplica o timeout
	ctx, cancel := context.WithTimeout(ctx, 45*time.Second) // Tempo maior para download
	defer cancel()

	// 2. Constrói o Payload para o whatsapp-web-ms (se usar POST)
	requestPayload := map[string]string{
		"messageId": modelRequest.MessageID,
		"number":    modelRequest.Number,
	}
	jsonPayload, _ := json.Marshal(requestPayload)

	// 3. Faz a requisição HTTP POST para o whatsapp-web-ms na rota de download
	// (Assumindo que voce criou uma rota /download-media)
	url := fmt.Sprintf("%s/download-media", c.WWEBJS_MS_URL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("falha na comunicacao com whatsapp-web-ms: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("whatsapp-web-ms retornou status de erro: %d", resp.StatusCode)
	}

	// 4. Le o corpo (os bytes do arquivo)
	audioBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler a midia do whatsapp-web-ms: %w", err)
	}

	return audioBytes, nil
}
