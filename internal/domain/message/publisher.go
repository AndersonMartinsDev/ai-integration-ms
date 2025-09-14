package message

import "context"

// MessageConsumer define a interface para consumir mensagens de uma fila.
type MessagePublisher interface {
	Publish(ctx context.Context, queueName string, data []byte) error
	Close() error
}
