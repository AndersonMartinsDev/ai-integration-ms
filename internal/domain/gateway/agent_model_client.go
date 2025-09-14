package gateway

import (
	"context"

	"ai-integration-ms/proto"

	"google.golang.org/grpc"
)

// WebhookProcessorClient define a interface para o serviço de processamento de webhooks.
// A camada de aplicação usará essa interface, sem saber que a implementação é gRPC.
type AgentClient interface {
	GetAgent(ctx context.Context, in *proto.GetAgentRequest, opts ...grpc.CallOption) (*proto.GetAgentResponse, error)
}
