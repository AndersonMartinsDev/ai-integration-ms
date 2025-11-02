package cmd

import (
	services "ai-integration-ms/internal/application/service"
	"ai-integration-ms/internal/infrastructure/ai/gemini"
	"ai-integration-ms/internal/infrastructure/rabbitmq"
	"ai-integration-ms/internal/infrastructure/repository"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MSCompose struct {
	RabbitMQURL     string
	MetaAccessToken string
}

func NewMSCompose() *MSCompose {
	return &MSCompose{
		RabbitMQURL:     os.Getenv("RABBITMQ_URL_FILE"),
		MetaAccessToken: os.Getenv("META_ACCESS_TOKEN"),
	}
}

func (manager MSCompose) GeminiServiceConfiguration(cacheService services.CacheService) *services.GeminiService {
	gemini_client := gemini.NewGeminiClient(gemini.GEMINI_MODEL)
	return services.NewGeminiService(gemini_client, cacheService)
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
