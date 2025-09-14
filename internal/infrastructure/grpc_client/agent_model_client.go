package grpc_client

import (
	"ai-integration-ms/internal/domain/gateway"
	"ai-integration-ms/proto"
	"context"

	"google.golang.org/grpc"
)

// WebhookClient implementa a interface domain/gateway.WebhookProcessorClient.
type AgentModelClient struct {
	client proto.AIAgentServiceClient
}

// NewWebhookClient cria uma nova instância do cliente gRPC.
func NewAgentModelClient(conn *grpc.ClientConn) gateway.AgentClient {
	return &AgentModelClient{
		client: proto.NewAIAgentServiceClient(conn),
	}
}

// GetAgent implements gateway.AgentModelClient.
func (a *AgentModelClient) GetAgent(ctx context.Context, in *proto.GetAgentRequest, opts ...grpc.CallOption) (*proto.GetAgentResponse, error) {
	return a.client.GetAgent(ctx, in, opts...)
}
