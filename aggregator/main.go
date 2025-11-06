package main

import (
	"context" // handle request contexts
	"fmt"
	"log"
	"net" // network operations

	aggregatorpb "aggregator/proto/aggregator/v1"
	"google.golang.org/grpc"
)

// server is used to implement the AggregatorServer interface
type AggregatorServer struct {
	aggregatorpb.UnimplementedAggregatorServer
}

// ReportResult implements the ReportResult RPC method
func (s *AggregatorServer) ReportResult(ctx context.Context, req *aggregatorpb.ReportResultRequest) (*aggregatorpb.ReportResultResponse, error) {
	// Log the incoming request
	log.Printf("Received ReportResult request:")
	log.Printf("  Task ID: %s", req.TaskId)
	log.Printf("  Attempt ID: %d", req.AttemptId)
	log.Printf("  Status: %s", req.Status)
	log.Printf("  Result Ref: %s", req.ResultRef)
	log.Printf("  Metrics: %s", req.Metrics)
	
	// Process the request (add your business logic here)
	// Return a response
	return &aggregatorpb.ReportResultResponse{
		Acknowledged: true,
		Message:      fmt.Sprintf("Task %s result received successfully", req.TaskId),
	}, nil
}

func main() {
	// Define the port to listen on
	port := ":50051"
	
	// Create a TCP listener
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()
	
	// Register our Aggregator service with the gRPC server
	aggregatorpb.RegisterAggregatorServer(grpcServer, &AggregatorServer{})
	
	log.Printf("gRPC server listening on %s", port)
	
	// Start serving (this blocks)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}