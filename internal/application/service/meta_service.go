package services

import (
	"ai-integration-ms/internal/domain/model"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
	// Assumindo que você tem acesso ao seu modelo e clientes
)

// Estrutura do MetaService (Com a configuracao do Meta)
type MetaService struct {
	MetaAccessToken string // Token de Acesso do Meta
	HttpClient      *http.Client
}

func NewMetaService(httpClient *http.Client, metaToken string) *MetaService {
	return &MetaService{
		MetaAccessToken: metaToken,
		HttpClient:      httpClient,
	}
}

// downloadAudioFile busca o arquivo de áudio da URL fornecida, aplicando autenticação e timeout.
func (s *MetaService) DownloadAudioFile(ctx context.Context, requestModel model.AiRequestModel) ([]byte, error) {
	// Definimos o timeout para a requisição de download (30 segundos)
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 1. Cria a Requisição HTTP GET
	req, err := http.NewRequestWithContext(ctx, "GET", requestModel.MediaUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar a requisição GET: %w", err)
	}

	// 2. Adiciona o Cabeçalho de Autenticação
	// Este passo é CRÍTICO para URLs de mídia do Meta.
	if s.MetaAccessToken == "" {
		slog.Error("AccessToken do Meta está vazio. Nao e possivel baixar a midia.")
		return nil, fmt.Errorf("token de acesso do Meta ausente")
	}
	req.Header.Add("Authorization", "Bearer "+s.MetaAccessToken)

	// Opcional: Adicionar um User-Agent.
	req.Header.Add("User-Agent", "ai-integration-ms-v1")

	// 3. Executa a Requisição
	resp, err := s.HttpClient.Do(req)
	if err != nil {
		// Inclui falhas de rede e timeout do contexto
		return nil, fmt.Errorf("erro na chamada HTTP de download: %w", err)
	}
	defer resp.Body.Close()

	// 4. Verifica o Status da Resposta
	if resp.StatusCode != http.StatusOK {
		// Tenta ler o corpo para logs detalhados de erro da API
		errorBody, _ := io.ReadAll(resp.Body)
		slog.Error("Download de audio falhou no Meta", "status", resp.StatusCode, "url", requestModel.MediaUrl, "response", string(errorBody))
		return nil, fmt.Errorf("Meta retornou erro %d ao tentar baixar a midia", resp.StatusCode)
	}

	// 5. Lê o Corpo da Resposta (Conteúdo do Arquivo)
	// Usamos io.ReadAll para carregar o conteúdo binário do arquivo na memória.
	audioBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler o corpo da midia: %w", err)
	}

	// Opcional: Verifique se o tamanho do arquivo não excede um limite razoável
	// if len(audioBytes) > MaxAudioFileSize { ... }

	return audioBytes, nil
}
