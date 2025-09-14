package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publisher é a implementação do publisher que se comunica com o RabbitMQ.
type Publisher struct {
	conn *amqp.Connection
}

// NewPublisher cria um novo publisher do RabbitMQ.
func NewPublisher(conn *amqp.Connection, amqpURL string) (*Publisher, error) {
	return &Publisher{conn: conn}, nil
}

// Publish envia os dados para a fila especificada.
func (p *Publisher) Publish(ctx context.Context, queueName string, data []byte) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return fmt.Errorf("falha ao abrir um canal: %w", err)
	}
	defer ch.Close()

	// Declara a fila (se ela já existir, nada acontece)
	q, err := ch.QueueDeclare(
		queueName, // nome da fila
		true,      // durável
		false,     // auto-delete
		false,     // exclusiva
		false,     // no-wait
		nil,       // argumentos
	)
	if err != nil {
		return fmt.Errorf("falha ao declarar a fila %s: %w", q.Name, err)
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		})
	if err != nil {
		return fmt.Errorf("falha ao publicar a mensagem: %w", err)
	}

	return nil
}

// Close fecha o canal e a conexão.
func (c *Publisher) Close() error {
	if c.conn != nil {
		c.conn.Close()
	}
	return nil
}
