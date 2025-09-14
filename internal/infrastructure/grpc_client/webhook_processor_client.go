package grpc_client

import (
	"ai-integration-ms/internal/domain/gateway"
	"ai-integration-ms/internal/domain/model"

	"google.golang.org/grpc"
)

// WebhookClient implementa a interface domain/gateway.WebhookProcessorClient.
type WebHookProcessorClientImpl struct {
	// client proto.AIAgentServiceClient
}

// NewWebhookClient cria uma nova instância do cliente gRPC.
func NewWebHookProcessorClientImpl(conn *grpc.ClientConn) gateway.WebHookProcessorClient {
	return &WebHookProcessorClientImpl{
		// client: proto.NewAIAgentServiceClient(conn),
	}
}

// GetSessionData implements gateway.WebHookProcessorClient.
func (w *WebHookProcessorClientImpl) GetSessionData(session_id string) *model.SessionModel {
	panic("unimplemented")
}

// // GetAgent implements gateway.AgentModelClient.
// func (a *AgentModelClient) GetAgent(ctx context.Context, in *proto.GetAgentRequest, opts ...grpc.CallOption) (*proto.GetAgentResponse, error) {
// 	return a.client.GetAgent(ctx, in, opts...)
// }
