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

type ServiceConfig interface {
	SetID(string)
	SetHost(string)
	SetPort(int)
	SetName(string)
}

/**
 * RegisterServiceWithConsul registers different types of services with Consul
 * 
 * TODO: Add credential implementation.
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

/**
 * GetServiceDetails retrieves details of a service registered in Consul.
 *
 * Parameters:
 *   - ctx: Context for the operation
 *   - serviceName: Name of the service to look up
 *   - log: Logger instance for recording operations
 *
 * Returns:
 *   - []*api.AgentService: List of service instances
 *   - error: Any error that occurred during the lookup
 */
func getServiceDetails(ctx context.Context, serviceName string, log logger.Zapper) ([]*api.AgentService, error) {
	// Create Consul client
	config := api.DefaultConfig()
	client, err := api.NewClient(config)
	if err != nil {
		log.Errorf(ctx, "Error creating Consul client: %v", err)
		return nil, fmt.Errorf("failed to create consul client: %w", err)
	}

	services, _, err := client.Health().Service(serviceName, "", true, &api.QueryOptions{})
	if err != nil {
		log.Errorf(ctx, "Error getting service details from Consul: %v", err)
		return nil, fmt.Errorf("failed to get service details: %w", err)
	}

	if len(services) == 0 {
		log.Warnf(ctx, "No healthy instances found for service: %s", serviceName)
		return nil, fmt.Errorf("no healthy instances found for service: %s", serviceName)
	}

	var serviceInstances []*api.AgentService
	for _, service := range services {
		serviceInstances = append(serviceInstances, service.Service)
		log.Infof(ctx, "Found service instance: ID=%s, Address=%s:%d",
			service.Service.ID,
			service.Service.Address,
			service.Service.Port,
		)
	}

	return serviceInstances, nil
}

/**
 * GetClientConfig is now generic and works with any config type that implements ServiceConfig
 */
func GetService[T ServiceConfig](ctx context.Context, client string, config T, log logger.Zapper) (T, error) {

	services, err := getServiceDetails(ctx, client, log)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to get client config or discover service: %w", err)
	}

	service := services[0]

	config.SetHost(service.Address)
	config.SetPort(service.Port)
	config.SetID(service.ID)
	config.SetName(service.Service)

	log.Infof(ctx, "Service config discovered from Consul: %s at %s:%d",
		client, service.Address, service.Port)

	return config, nil
}

/**
 * DeregisterService removes a service from Consul registration
 */
func DeregisterService(ctx context.Context, serviceID string, log logger.Zapper) error {
	config := api.DefaultConfig()
	client, err := api.NewClient(config)
	if err != nil {
		log.Errorf(ctx, "Error creating Consul client: %v", err)
		return fmt.Errorf("failed to create consul client: %w", err)
	}

	err = client.Agent().ServiceDeregister(serviceID)
	if err != nil {
		log.Errorf(ctx, "Failed to deregister service from Consul: %v", err)
		return fmt.Errorf("failed to deregister service: %w", err)
	}

	log.Infof(ctx, "Service deregistered from Consul: %s", serviceID)
	return nil
}
