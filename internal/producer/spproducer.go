package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/snnus/mainservice/config"
)

type TicketMessage struct {
	Ticket       string `json:"ticket"`
	OfficeNumber string `json:"officeNumber"`
	Timestamp    string `json:"timestamp"`
}

type SPProducer struct {
	writer *kafka.Writer
}

func NewSPProducer(cfg *config.Config) *SPProducer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Kafka.Broker),
		Topic:        cfg.Kafka.Topic,
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    cfg.Kafka.BatchSize,
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: kafka.RequireOne,
		Async:        false,
	}
	return &SPProducer{
		writer: writer,
	}
}

func (kp *SPProducer) PublishTicket(ctx context.Context, ticket, officeNumber string) error {
	msg := TicketMessage{
		Ticket:       ticket,
		OfficeNumber: officeNumber,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
	}
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = kp.writer.WriteMessages(ctx, kafka.Message{
		Value: jsonData,
	})

	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

func (kp *SPProducer) Close() error {
	return kp.writer.Close()
}
