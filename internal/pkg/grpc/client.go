package grpc

import (
	"context"
	"fmt"
	"pkg/discovery"
	"pkg/logger"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcClient struct {
	conn *grpc.ClientConn
}

type GrpcClientConfig map[string]GrpcConfig

type GrpcClient interface {
	GetGrpcConnection() *grpc.ClientConn
	Close() error
}

/**
 * NewGrpcClient to Call Grpc Server
 */
func NewGrpcClient(ctx context.Context, client string, clientsConfig *GrpcClientConfig, log logger.Zapper) (GrpcClient, error) {
	var config GrpcConfig
	var ok bool

	if config, ok = (*clientsConfig)[client]; !ok {

		discoveredConfig, err := discovery.GetService(ctx, client, &GrpcConfig{}, log)

		if err != nil {
			return nil, err
		}

		config = *discoveredConfig
		(*clientsConfig)[client] = config
	}

	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", config.Host, config.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)

	if err != nil {
		log.Errorf(ctx, "GRPC failed connecting grpc server %v", err)
		return nil, err
	}

	return &grpcClient{conn: conn}, nil
}

func (g *grpcClient) GetGrpcConnection() *grpc.ClientConn {
	return g.conn
}

func (g *grpcClient) Close() error {
	return g.conn.Close()
}
