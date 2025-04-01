package kafka

import (
	"context"
	"fmt"
	"pkg/logger"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
	log    logger.Zapper
	cfg    KafkaConfig
}

type KafkaConsumer struct {
	reader *kafka.Reader
	log    logger.Zapper
	cfg    KafkaConfig
}

// NewKafkaProducer creates a new Kafka producer
func NewKafkaProducer(cfg KafkaConfig, log logger.Zapper) *KafkaProducer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      cfg.Brokers,
		Topic:        cfg.Topic,
		BatchSize:    100,
		BatchTimeout: 1 * time.Second,
		Async:        true,
	})

	return &KafkaProducer{
		writer: writer,
		log:    log,
		cfg:    cfg,
	}
}

// NewKafkaConsumer creates a new Kafka consumer
func NewKafkaConsumer(cfg KafkaConfig, log logger.Zapper) *KafkaConsumer {

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     cfg.Brokers,
		Topic:       cfg.Topic,
		GroupID:     cfg.GroupID,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		MaxWait:     1 * time.Second,
		StartOffset: kafka.FirstOffset,
	})

	return &KafkaConsumer{
		reader: reader,
		log:    log,
		cfg:    cfg,
	}
}

// Produce sends a message to Kafka
func (p *KafkaProducer) Produce(ctx context.Context, key, value []byte) error {
	msg := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}

	err := p.writer.WriteMessages(ctx, msg)
	if err != nil {
		p.log.Errorf(ctx, "Failed to write message to Kafka: %v", err)
		return fmt.Errorf("failed to write message: %w", err)
	}

	p.log.Infof(ctx, "Successfully produced message to topic %s", p.cfg.Topic)
	return nil
}

// ProduceBatch sends multiple messages to Kafka
func (p *KafkaProducer) ProduceBatch(ctx context.Context, messages []kafka.Message) error {
	err := p.writer.WriteMessages(ctx, messages...)
	if err != nil {
		p.log.Errorf(ctx, "Failed to write batch messages to Kafka: %v", err)
		return fmt.Errorf("failed to write batch messages: %w", err)
	}

	p.log.Infof(ctx, "Successfully produced %d messages to topic %s", len(messages), p.cfg.Topic)
	return nil
}

// Consume reads messages from Kafka
// TODO: Move to go routienes.
func (c *KafkaConsumer) Consume(ctx context.Context, handler func([]byte) error) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				return nil
			}
			c.log.Errorf(ctx, "Failed to read message from Kafka: %v", err)
			return fmt.Errorf("failed to read message: %w", err)
		}

		if err := handler(msg.Value); err != nil {
			c.log.Errorf(ctx, "Failed to process message: %v", err)
			continue
		}

		c.log.Infof(ctx, "Successfully processed message from topic %s", c.cfg.Topic)
	}
}

// ConsumeWithCommit reads messages and commits offsets
// TODO: Move to separate go routiene.
func (c *KafkaConsumer) ConsumeWithCommit(ctx context.Context, handler func([]byte) error) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				return nil
			}
			c.log.Errorf(ctx, "Failed to read message from Kafka: %v", err)
			return fmt.Errorf("failed to read message: %w", err)
		}

		if err := handler(msg.Value); err != nil {
			c.log.Errorf(ctx, "Failed to process message: %v", err)
			continue
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			c.log.Errorf(ctx, "Failed to commit message: %v", err)
			continue
		}

		c.log.Infof(ctx, "Successfully processed and committed message from topic %s", c.cfg.Topic)
	}
}

// Close closes the Kafka connection
func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}
