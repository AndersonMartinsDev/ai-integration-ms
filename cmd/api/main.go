package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync" // <-- Adicione essa importação
	"syscall"

	"ai-integration-ms/cmd"
	"ai-integration-ms/internal/infrastructure/configuration"
	"ai-integration-ms/internal/infrastructure/database"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// 1. Configuração e inicialização de todas as dependências
	// configuration.LoadEnv()
	configuration.LoadLogger()
	configuration.LoadRedis()
	configuration.ConfigGenerativeGemini()

	// Use WaitGroup para garantir que todas as goroutines encerrem
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	ms_compose := cmd.NewMSCompose()

	conn, err := amqp.Dial(ms_compose.RabbitMQURL)
	if err != nil {
		log.Fatalf("Falha ao conectar no RabbitMQ: %v", err)
	}

	cacheService := ms_compose.CacheServiceConfiguration()
	geminiService := ms_compose.GeminiServiceConfiguration(*cacheService)
	message_processor := ms_compose.MessageProcessorConfiguration(conn, *geminiService)

	// Inicia a goroutine do consumidor e a adiciona ao WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		message_processor.ProcessWhatsAppMessages(ctx)
	}()

	// 6. BLOQUEAR a execução da main até que um sinal seja recebido
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// 7. Lógica de encerramento
	slog.Info("Sinal de interrupção recebido, iniciando o encerramento...")

	// 7.1. Cancela o contexto para sinalizar às goroutines para pararem
	cancel()

	// 7.2. Espera todas as goroutines (o consumidor) finalizarem
	wg.Wait()

	// 7.3. Fecha os recursos de forma limpa, APÓS as goroutines terminarem
	message_processor.Publisher.Close()
	message_processor.Consumer.Close()
	conn.Close()
	database.CACHE.Close()

	slog.Info("Serviço encerrado com sucesso!")
}
