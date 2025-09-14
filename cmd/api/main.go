package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"ai-integration-ms/cmd"
	"ai-integration-ms/internal/infrastructure/configuration"
	"ai-integration-ms/internal/infrastructure/database"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// 1. Configuração e inicialização de todas as dependências
	configuration.LoadEnv()
	configuration.LoadLogger()
	configuration.LoadRedis()
	configuration.ConfigGenerativeGemini()
	ctx, cancel := context.WithCancel(context.Background())

	ms_compose := cmd.NewMSCompose()

	conn, err := amqp.Dial(ms_compose.RabbitMQURL)
	if err != nil {
		log.Fatalf("Falha ao conectar no RabbitMQ: %v", err)
	}

	cacheService := ms_compose.CacheServiceConfiguration()
	geminiService := ms_compose.GeminiServiceConfiguration(*cacheService)
	message_processor := ms_compose.MessageProcessorConfiguration(conn, *geminiService)

	go message_processor.ProcessWhatsAppMessages(ctx)
	defer message_processor.Publisher.Close()
	defer message_processor.Consumer.Close()
	defer conn.Close()
	defer database.CACHE.Close()

	defer cancel()

	// 6. BLOQUEAR a execução da main até que um sinal seja recebido
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan // Esta linha bloqueia a thread principal

	// 7. Lógica de encerramento
	slog.Info("Sinal de interrupção recebido, encerrando o serviço...")
}
