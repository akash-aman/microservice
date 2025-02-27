package discovery

import (
	"context"
	"fmt"
	"pkg/logger"

	"github.com/hashicorp/consul/api"
)

/**
 * ServiceType defines the type of service being registered
 */
type ServiceType string

const (
	HTTPService  ServiceType = "http"
	GRPCService  ServiceType = "grpc"
	KafkaService ServiceType = "kafka"
	TCPService   ServiceType = "tcp"
	RabbitMQ     ServiceType = "rabbitmq"
	RedisService ServiceType = "redis"
)

/**
 * RegisterServiceWithConsul registers different types of services with Consul
 */
func RegisterServiceWithConsul(ctx context.Context, serviceName, serviceID, address string, port int, serviceType ServiceType, log logger.Zapper) error {
	config := api.DefaultConfig()
	client, err := api.NewClient(config)
	if err != nil {
		log.Errorf(ctx, "Error creating Consul client: %v", err)
		return err
	}

	/**
	 * Define the service registration
	 */
	reg := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: address,
		Port:    port,
	}

	/**
	 * Register the service with Consul.
	 */
	switch serviceType {
	case HTTPService:
		reg.Check = &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("%s:%d/health", address, port),
			Interval: "10s",
			Timeout:  "5s",
		}
	case GRPCService:
		reg.Check = &api.AgentServiceCheck{
			GRPC:       fmt.Sprintf("%s:%d", address, port),
			GRPCUseTLS: false,
			Interval:   "10s",
			Timeout:    "5s",
		}
	case KafkaService, TCPService, RabbitMQ, RedisService:
		reg.Check = &api.AgentServiceCheck{
			TCP:      fmt.Sprintf("%s:%d", address, port),
			Interval: "10s",
			Timeout:  "5s",
		}
	}

	/**
	 * Register the service with Consul
	 */
	err = client.Agent().ServiceRegister(reg)
	if err != nil {
		log.Errorf(ctx, "Failed to register service with Consul: %v", err)
		return err
	}

	log.Infof(ctx, "Service registered with Consul: %s (%s) on %s:%d", serviceName, serviceType, address, port)
	return nil
}
