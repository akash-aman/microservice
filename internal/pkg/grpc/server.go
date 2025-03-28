package grpc

/**
 * https://signoz.io/blog/opentelemetry-grpc-golang/
 *
 */

import (
	"context"
	"fmt"
	"net"
	"pkg/discovery"
	"pkg/helper"
	"pkg/logger"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

const (
	maxConnectionIdle = 5
	gRPCTimeout       = 15
	maxConnectionAge  = 5
	gRPCTime          = 10
)

type GrpcConfig struct {
	Port        int    `mapstructure:"port"`
	Host        string `mapstructure:"host"`
	Development bool   `mapstructure:"development"`
}

// Close implements GrpcClient.
func (c GrpcConfig) Close() error {
	panic("unimplemented")
}

// GetGrpcConnection implements GrpcClient.
func (c GrpcConfig) GetGrpcConnection() *grpc.ClientConn {
	panic("unimplemented")
}

func (c *GrpcConfig) SetHost(host string) {
	c.Host = host
}

func (c *GrpcConfig) SetPort(port int) {
	c.Port = port
}

func (c *GrpcConfig) SetID(id string) {
	panic("unimplemented")
}

func (c *GrpcConfig) SetName(name string) {
	panic("unimplemented")
}

type GrpcServer struct {
	Grpc   *grpc.Server
	Config *GrpcConfig
	Log    logger.Zapper
}

func NewGrpcServer(log logger.Zapper, config *GrpcConfig) *GrpcServer {

	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				MaxConnectionIdle: maxConnectionIdle * time.Minute,
				Timeout:           gRPCTimeout * time.Second,
				MaxConnectionAge:  maxConnectionAge * time.Minute,
				Time:              gRPCTime * time.Minute,
			},
		),
	)

	return &GrpcServer{Grpc: s, Config: config, Log: log}
}

func (s *GrpcServer) RunGrpcServer(ctx context.Context, log logger.Zapper, service string, configGrpc ...func(grpcServer *grpc.Server)) error {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port))
	if err != nil {
		return errors.Wrap(err, "net.Listen")
	}

	if len(configGrpc) > 0 {
		grpcFunc := configGrpc[0]
		if grpcFunc != nil {
			grpcFunc(s.Grpc)
		}
	}

	if s.Config.Development {
		reflection.Register(s.Grpc)
	}

	if len(configGrpc) > 0 {
		grpcFunc := configGrpc[0]
		if grpcFunc != nil {
			grpcFunc(s.Grpc)
		}
	}

	go func() {
		for range ctx.Done() {
			s.Log.Infof(ctx, "shutting down grpc PORT: {%s}", s.Config.Port)
			s.shutdown(service)
			s.Log.Info(ctx, "grpc exited properly")
			return
		}
	}()

	s.Log.Infof(ctx, "grpc server is listening on port: %d", s.Config.Port)

	err = s.Grpc.Serve(listen)

	if err != nil {
		s.Log.Error(ctx, fmt.Sprintf("[grpcServer_RunGrpcServer.Serve] grpc server serve error: %+v", err))
	}

	if err := discovery.RegisterServiceWithConsul(
		ctx, service,
		fmt.Sprintf("%s-%s", service, helper.GetMachineID()),
		fmt.Sprintf("http://%s", s.Config.Host),
		s.Config.Port,
		discovery.HTTPService,
		log,
	); err != nil {
		log.Errorf(ctx, "Error registering with Consul: %v", err)
	}

	return err
}

func (s *GrpcServer) shutdown(service string) {

	if err := discovery.DeregisterService(context.Background(), service, s.Log); err != nil {
		s.Log.Errorf(context.Background(), "Failed to deregister service: %v", err)
	}

	s.Grpc.GracefulStop()
}
