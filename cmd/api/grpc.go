package main

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func setupGRPC(cfg *Config) (*grpc.ClientConn, error) {
	dialOps := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	}
	conn, err := grpc.NewClient(cfg.Service.ProductAddr, dialOps...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
