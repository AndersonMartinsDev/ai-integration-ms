package cmd

import (
	services "ai-integration-ms/internal/application/service"
	"ai-integration-ms/internal/infrastructure/ai/gemini"
	"ai-integration-ms/internal/infrastructure/configuration"
	"ai-integration-ms/internal/infrastructure/grpc_client"
	"ai-integration-ms/internal/infrastructure/rabbitmq"
	"ai-integration-ms/internal/infrastructure/repository"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MSCompose struct {
	AgentModelMSUrl string
	RabbitMQURL     string
	MetaAccessToken string
}

func NewMSCompose() *MSCompose {
	return &MSCompose{
		AgentModelMSUrl: os.Getenv("AGENT_MODEL_URL"),
		RabbitMQURL:     configuration.GetSecret("RABBITMQ_URL_FILE"),
		MetaAccessToken: os.Getenv("META_ACCESS_TOKEN"),
	}
}

func (manager MSCompose) GeminiServiceConfiguration(cacheService services.CacheService) *services.GeminiService {
	agent_model_client := grpc_client.NewAgentModelClient(grpc_client.GrcpConnection(manager.AgentModelMSUrl))
	gemini_client := gemini.NewGeminiClient(gemini.GEMINI_MODEL)
	return services.NewGeminiService(agent_model_client, gemini_client, cacheService)
}

func (manager MSCompose) CacheServiceConfiguration() *services.CacheService {
	redis_repository := repository.NewRedisRepository()
	return services.NewCacheService(redis_repository)
}

func (manager MSCompose) MessageProcessorConfiguration(conn *amqp.Connection, geminiService services.GeminiService) *services.MessageProcessor {
	// httpClient := &http.Client{Timeout: 60 * time.Second}
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
