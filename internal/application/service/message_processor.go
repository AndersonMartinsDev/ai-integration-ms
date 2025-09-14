package services

import (
	"ai-integration-ms/internal/domain/message"
	"context"
	"encoding/json"
	"log"
	"log/slog"
)

type MessageProcessor struct {
	GeminiService GeminiService
	Consumer      message.MessageConsumer
	Publisher     message.MessagePublisher
}

func NewMessageProcessor(service GeminiService, consumer message.MessageConsumer, publisher message.MessagePublisher) *MessageProcessor {
	return &MessageProcessor{
		GeminiService: service,
		Consumer:      consumer,
		Publisher:     publisher,
	}
}

func (process *MessageProcessor) ProcessWhatsAppMessages(ctx context.Context) {
	msgs, err := process.Consumer.Consume(ctx)
	if err != nil {
		log.Fatalf("Falha ao iniciar o consumo de mensagens: %v", err)
	}

	slog.Info("Serviço Gemini Integration iniciado, esperando por mensagens...")

	for msgPayload := range msgs {
		slog.Info("Mensagem recebida, desserializando...")

		var inputMsg message.InputMessage
		err := json.Unmarshal(msgPayload, &inputMsg)
		if err != nil {
			slog.Error("Falha ao desserializar a mensagem", "error", err)
			continue
		}

		slog.Info("Mensagem desserializada, processando com o Gemini...")

		response, err := process.GeminiService.GenerateText(
			inputMsg.AgentID,
			inputMsg.SessionKey,
			inputMsg.FontNumber,
			inputMsg.Message,
		)

		if err != nil {
			slog.Error("Falha ao processar mensagem", "error", err)
		}

		var outputMessage message.OutPutMessage
		outputMessage.Message = response
		outputMessage.PhoneNumber = inputMsg.SessionKey
		outputMessage.FontNumber = inputMsg.FontNumber
		outputMessage.MessageType = inputMsg.MessageType

		publisherMessage, erro := json.Marshal(outputMessage)
		if erro != nil {
			slog.Error("Falha ao processar mensagem", "error", err)
		}
		process.Publisher.Publish(ctx, "whatsapp-generated-message", publisherMessage)

	}
}
