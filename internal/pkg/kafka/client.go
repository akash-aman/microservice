package kafka

import (
	"context"
	"fmt"
	"pkg/logger"
	"sync"

	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	producers map[string]*KafkaProducer
	consumers map[string]*KafkaConsumer
	mu        sync.RWMutex
	log       logger.Zapper
}

func NewKafkaClient(log logger.Zapper) *KafkaClient {
	return &KafkaClient{
		producers: make(map[string]*KafkaProducer),
		consumers: make(map[string]*KafkaConsumer),
		log:       log,
	}
}

// getControllerBroker finds the controller broker in the cluster
func (c *KafkaClient) getControllerBroker(ctx context.Context, brokers []string) (string, error) {
	for _, broker := range brokers {
		conn, err := kafka.Dial("tcp", broker)
		if err != nil {
			c.log.Errorf(ctx, "Failed to connect to broker %s: %v", broker, err)
			continue
		}
		defer conn.Close()

		controller, err := conn.Controller()
		if err != nil {
			c.log.Errorf(ctx, "Failed to get controller from broker %s: %v", broker, err)
			continue
		}

		controllerAddr := fmt.Sprintf("%s:%d", controller.Host, controller.Port)
		c.log.Infof(ctx, "Found controller broker: %s", controllerAddr)
		return controllerAddr, nil
	}

	return "", fmt.Errorf("failed to find controller broker")
}

func (c *KafkaClient) CreateTopicIfNotExists(ctx context.Context, cfg KafkaConfig) error {
	controllerAddr, err := c.getControllerBroker(ctx, cfg.Brokers)
	if err != nil {
		return fmt.Errorf("failed to get controller broker: %w", err)
	}

	conn, err := kafka.Dial("tcp", controllerAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to controller: %w", err)
	}
	defer conn.Close()

	// Check if topic exists
	partitions, err := conn.ReadPartitions(cfg.Topic)
	if err == nil && len(partitions) > 0 {
		c.log.Infof(ctx, "Topic %s already exists", cfg.Topic)
		return nil
	}

	// Create topic configuration
	topic := kafka.TopicConfig{
		Topic:             cfg.Topic,
		NumPartitions:     cfg.PartitionCount,
		ReplicationFactor: cfg.ReplicationFactor,
		ConfigEntries: []kafka.ConfigEntry{
			{
				ConfigName:  "retention.ms",
				ConfigValue: fmt.Sprintf("%d", cfg.RetentionHours*60*60*1000),
			},
		},
	}

	err = conn.CreateTopics(topic)
	if err != nil {
		return fmt.Errorf("failed to create topic: %w", err)
	}

	c.log.Infof(ctx, "Created topic %s with %d partitions and replication factor %d",
		cfg.Topic, cfg.PartitionCount, cfg.ReplicationFactor)
	return nil
}

func (c *KafkaClient) GetProducer(ctx context.Context, topic string, cfg KafkaConfig) (*KafkaProducer, error) {
	if err := c.CreateTopicIfNotExists(ctx, cfg); err != nil {
		return nil, err
	}

	c.mu.RLock()
	if producer, exists := c.producers[topic]; exists {
		c.mu.RUnlock()
		return producer, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	if producer, exists := c.producers[topic]; exists {
		return producer, nil
	}

	producer := NewKafkaProducer(cfg, c.log)
	c.producers[topic] = producer
	return producer, nil
}

func (c *KafkaClient) GetConsumer(ctx context.Context, topic string, cfg KafkaConfig) (*KafkaConsumer, error) {

	if err := c.CreateTopicIfNotExists(ctx, cfg); err != nil {
		return nil, err
	}

	c.mu.RLock()
	if consumer, exists := c.consumers[topic]; exists {
		c.mu.RUnlock()
		return consumer, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double check after acquiring write lock
	if consumer, exists := c.consumers[topic]; exists {
		return consumer, nil
	}

	consumer := NewKafkaConsumer(cfg, c.log)
	c.consumers[topic] = consumer
	return consumer, nil
}

func (c *KafkaClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var errs []error

	// Close all producers
	for topic, producer := range c.producers {
		if err := producer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close producer for topic %s: %w", topic, err))
		}
	}

	// Close all consumers
	for topic, consumer := range c.consumers {
		if err := consumer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close consumer for topic %s: %w", topic, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors while closing Kafka clients: %v", errs)
	}

	return nil
}
