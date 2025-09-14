package cmd

import (
	services "ai-integration-ms/internal/application/service"
	"ai-integration-ms/internal/infrastructure/ai/gemini"
	"ai-integration-ms/internal/infrastructure/grpc_client"
	"ai-integration-ms/internal/infrastructure/rabbitmq"
	"ai-integration-ms/internal/infrastructure/repository"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MSCompose struct {
	AgentModelMSUrl       string
	WebHookProcessorMSURL string
	RabbitMQURL           string
}

func NewMSCompose() *MSCompose {
	return &MSCompose{
		AgentModelMSUrl:       os.Getenv("AGENT_MODEL_URL"),
		WebHookProcessorMSURL: os.Getenv("WEBHOOK_PROCESSOR_URL"),
		RabbitMQURL:           os.Getenv("RABBITMQ_URL"),
	}
}

func (manager MSCompose) GeminiServiceConfiguration(cacheService services.CacheService) *services.GeminiService {
	gemini_client := gemini.NewGeminiClient(gemini.GEMINI_MODEL)
	agent_model_client := grpc_client.NewAgentModelClient(grpc_client.GrcpConnection(manager.AgentModelMSUrl))
	webhook_processor_client := grpc_client.NewWebHookProcessorClientImpl(grpc_client.GrcpConnection(manager.WebHookProcessorMSURL))
	return services.NewGeminiService(gemini_client, agent_model_client, cacheService, webhook_processor_client)
}

func (manager MSCompose) CacheServiceConfiguration() *services.CacheService {
	redis_repository := repository.NewRedisRepository()
	return services.NewCacheService(redis_repository)
}

func (manager MSCompose) MessageProcessorConfiguration(conn *amqp.Connection, geminiService services.GeminiService) *services.MessageProcessor {
	publisher, err := rabbitmq.NewPublisher(conn, "whatsapp-outbound-messages")
	if err != nil {
		log.Fatalf("Falha ao criar publicador RabbitMQ: %v", err)
	}
	consumer, err := rabbitmq.NewConsumer(conn, "webhook-whatsapp-messages")
	if err != nil {
		log.Fatalf("Falha ao criar consumidor RabbitMQ: %v", err)
	}
	return services.NewMessageProcessor(geminiService, consumer, publisher)
}
