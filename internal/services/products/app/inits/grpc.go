package inits

import (
	"pkg/grpc"
	"pkg/logger"
	"products/app/grpc/server"
	"products/app/grpc/server/proto"
	"products/conf"
)

func ConfigGrpcServer(cfg *conf.Config, log logger.Zapper, grpcServer *grpc.GrpcServer) {
	productGrpcService := server.NewProductGrpcServerService(log, *cfg)

	// register the server.
	proto.RegisterProductServiceServer(grpcServer.Grpc, productGrpcService)
}
