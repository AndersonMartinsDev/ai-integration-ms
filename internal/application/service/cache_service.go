package services

import (
	"ai-integration-ms/internal/domain/model"
	"ai-integration-ms/internal/infrastructure/repository"
	"encoding/json"
	"fmt"

	"github.com/google/generative-ai-go/genai"
)

type CacheService struct {
	Repository repository.RedisRepository
}

func NewCacheService(repository repository.RedisRepository) *CacheService {
	return &CacheService{
		Repository: repository,
	}
}

func (cache *CacheService) SaveHistory(session model.SessionModel, history []*genai.Content) error {
	historyJSON, err := json.Marshal(history)
	if err != nil {
		return fmt.Errorf("erro ao serializar histórico genai para json.RawMessage: %w", err)
	}

	session.History = historyJSON

	if _, err := cache.Repository.Save(&session); err != nil {
		return fmt.Errorf("erro ao salvar sessão e histórico: %w", err)
	}
	return nil
}

func (cache *CacheService) GetHistory(fromNumber, fontNumber string) (model.SessionModel, []*genai.Content, error) {
	keySession := fmt.Sprintf("%s,%s", fromNumber, fontNumber)
	register, _ := cache.Repository.Get(keySession)
	if register == nil {
		return model.SessionModel{
			KeyId:       keySession,
			PhoneNumber: fontNumber,
		}, []*genai.Content{}, nil
	}

	var record []*struct {
		Parts []interface{}
		Role  string
	}

	err := json.Unmarshal(register.History, &record)
	if err != nil {
		return model.SessionModel{}, []*genai.Content{}, err
	}

	var history []*genai.Content

	for _, content := range record {

		hist := &genai.Content{
			Role:  content.Role,
			Parts: []genai.Part{},
		}

		for _, text := range content.Parts {
			hist.Parts = append(hist.Parts, genai.Text(text.(string)))
		}
		history = append(history, hist)
	}
	return model.SessionModel{
		KeyId:   keySession,
		AgentId: register.AgentId,
	}, history, nil
}
