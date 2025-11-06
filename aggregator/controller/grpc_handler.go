package controller

import (
	"google.golang.org/grpc"
	aggregatorpb "aggregator/proto/aggregator/v1"
)

func RegisterGRPCServers(grpcServer *grpc.Server) {
	// Register our Aggregator service with the gRPC server
	aggregatorpb.RegisterAggregatorServer(grpcServer, &AggregatorServer{})

	//Register more services here ->
}