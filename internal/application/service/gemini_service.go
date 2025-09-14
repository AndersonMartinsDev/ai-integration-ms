package services

import (
	"ai-integration-ms/internal/domain/gateway"
	"ai-integration-ms/internal/domain/model"
	"ai-integration-ms/proto"
	"context"
	"fmt"
	"log/slog"
	"sync"
)

type GeminiService struct {
	AgentClient            gateway.AgentClient
	GeminiClient           gateway.GeminiClient
	WebHookProcessorClient gateway.WebHookProcessorClient
	CacheService           CacheService
}

func NewGeminiService(
	geminiClient gateway.GeminiClient,
	agentClient gateway.AgentClient,
	cache CacheService,
	webhookProcessorClient gateway.WebHookProcessorClient) *GeminiService {
	return &GeminiService{
		AgentClient:            agentClient,
		GeminiClient:           geminiClient,
		WebHookProcessorClient: webhookProcessorClient,
		CacheService:           cache,
	}
}

func (s *GeminiService) GenerateText(agentId uint, keyId, fontNumber, prompt string) (string, error) {

	var wg sync.WaitGroup
	wg.Add(1)

	context := context.Background()

	session, history, err := s.CacheService.GetHistory(keyId, fontNumber)
	if err != nil {
		slog.Error("Não foi possível recuperar historico de conversa!")
	}

	var instructions string
	go func() {
		defer wg.Done()
		instructions, err = s.getAgent(agentId, session.PhoneNumber)
		if err != nil {
			slog.Error("Erro ao montar instruções!")
		}
	}()

	wg.Wait()

	response, updatedHistory, err := s.GeminiClient.GenerateContent(context, instructions, history, prompt)

	if err != nil {
		slog.Error(fmt.Sprintf("Há um erro ao tentar gerar sua mensagem de resposta com a IA :  %s", err))
	}

	go func() {
		session.AgentId = agentId
		s.CacheService.SaveHistory(session, updatedHistory)
	}()

	return response, nil
}

func (s GeminiService) getAgent(agentId uint, uuid_user string) (string, error) {
	context := context.Background()
	req := &proto.GetAgentRequest{
		Id:           uint64(agentId),
		UuidUser:     uuid_user,
		Instructions: true,
	}

	result, err := s.AgentClient.GetAgent(context, req)
	if err != nil {
		slog.Error(err.Error())
		slog.Error("Não foi encontrado nenhum agente!")
		return "", err
	}

	// agent := result.Agent
	agent := &model.AIAgent{
		Id:                 uint64(agentId),
		Name:               result.GetAgent().Name,
		CompanyName:        result.GetAgent().CompanyName,
		CompanyDescription: result.GetAgent().CompanyDescription,
		BehaviourIa:        result.GetAgent().BehaviourIa,
		CompanyUrl:         result.GetAgent().CompanyUrl,
		Instructions:       result.GetAgent().Instructions,
	}
	return agent.GetSystemInstructions(), nil
}
