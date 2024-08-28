package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func setupGRPC(cfg *Config) (*grpc.ClientConn, error) {
	dialOps := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.NewClient(cfg.Service.ProductAddr, dialOps...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
