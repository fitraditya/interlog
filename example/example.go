package main

import (
	"net"

	"golang.org/x/net/context"

	"github.com/fitraditya/interlog"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	grpc_health "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type API struct{}

func (a *API) Check(ctx context.Context, in *grpc_health.HealthCheckRequest) (*grpc_health.HealthCheckResponse, error) {
	return &grpc_health.HealthCheckResponse{Status: grpc_health.HealthCheckResponse_SERVING}, nil
}

func (a *API) Watch(in *grpc_health.HealthCheckRequest, _ grpc_health.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "Unimplemented")
}

func main() {
	net, _ := net.Listen("tcp", ":9000")
	log := interlog.New()
	opt := []grpc.ServerOption{
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(log.Unary()),
		),
	}

	srv := grpc.NewServer(opt...)

	grpc_health.RegisterHealthServer(srv, &API{})

	_ = srv.Serve(net)
}
