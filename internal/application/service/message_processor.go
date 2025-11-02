package services

import (
	"ai-integration-ms/internal/domain/message"
	"ai-integration-ms/internal/domain/model"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"
)

type MessageProcessor struct {
	AIService     FileServiceInterface
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
	for {
		select {
		case <-ctx.Done():
			slog.Info("Contexto cancelado, encerrando o loop de consumo.")
			return
		default:
			err := process.processWhatsAppMessages(ctx)
			if err != nil {
				slog.Error("Falha na conexão ou consumo, tentando reconectar em 5 segundos...", "error", err)
				time.Sleep(5 * time.Second)
			} else {
				// Se a função retornar sem erro (o que é improvável em um loop de consumo),
				// a gente dá um pequeno tempo para evitar um loop muito rápido.
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func (process *MessageProcessor) processWhatsAppMessages(ctx context.Context) error {
	msgs, err := process.Consumer.Consume(ctx)
	if err != nil {
		log.Fatalf("Falha ao iniciar o consumo de mensagens: %v", err)
		return err
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

		messageInput := inputMsg.Message

		if inputMsg.MessageType == "audio" {
			service, err := process.featureDetection(inputMsg.OriginMessage)
			if err != nil {
				slog.Error("Falha ao recuperar feature para baixar audio: ", err)
				return err
			}
			audio, err := service.DownloadAudioFile(ctx, model.AiRequestModel{
				MediaUrl:  inputMsg.MediaUrl,
				MessageID: inputMsg.MessageID,
				Number:    inputMsg.FontNumber,
			})

			if err != nil {
				slog.Error("Falha ao baixar arquivo de audio", inputMsg.Message, inputMsg.Message, "error", err)
				return err
			}
			audioText, err := process.GeminiService.TranscribeAudio(ctx, audio)
			if err != nil {
				slog.Error("Falha a transcrever arquivo de audio", inputMsg.Message, inputMsg.Message, "error", err)
				return err
			}
			messageInput = audioText
		}

		response, err := process.GeminiService.GenerateText(
			inputMsg.AgentID,
			inputMsg.SessionKey,
			inputMsg.FontNumber,
			messageInput,
		)
		if err != nil {
			slog.Error("Falha ao processar mensagem", "error", err)
		}

		var outputMessage message.OutPutMessage
		outputMessage.Message = response
		outputMessage.PhoneNumber = inputMsg.SessionKey
		outputMessage.FontNumber = inputMsg.FontNumber
		outputMessage.MessageType = inputMsg.MessageType
		outputMessage.OriginMessage = inputMsg.OriginMessage

		publisherMessage, erro := json.Marshal(outputMessage)
		if erro != nil {
			slog.Error("Falha ao processar mensagem", "error", err)
		}
		process.Publisher.Publish(ctx, "whatsapp-generated-message", publisherMessage)

	}
	slog.Info("Consumidor de webhook encerrado.")
	// Retorna um erro para o loop principal, indicando que o consumo terminou.
	// Isso sinaliza que o canal foi fechado e que uma reconexão é necessária.
	return fmt.Errorf("canal de consumo fechado")
}

func (process *MessageProcessor) featureDetection(origin string) (FileServiceInterface, error) {
	httpClient := &http.Client{Timeout: 60 * time.Second}
	switch origin {
	case "META":
		return NewMetaService(httpClient, "PRECISA SER FEITA A LOGICA AINDA"), nil
	case "WWEBJS":
		return NewWWEBJSService(httpClient), nil
	default:
		return nil, fmt.Errorf("origem desconhecida: %s", origin)
	}
}
