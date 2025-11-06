package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "aggregator/proto/aggregator/v1"
	"google.golang.org/grpc"
)

// server is used to implement the AggregatorServer interface
type server struct {
	pb.UnimplementedAggregatorServer
}

// ReportResult implements the ReportResult RPC method
func (s *server) ReportResult(ctx context.Context, req *pb.ReportResultRequest) (*pb.ReportResultResponse, error) {
	// Log the incoming request
	log.Printf("Received ReportResult request:")
	log.Printf("  Task ID: %s", req.TaskId)
	log.Printf("  Attempt ID: %d", req.AttemptId)
	log.Printf("  Status: %s", req.Status)
	log.Printf("  Result Ref: %s", req.ResultRef)
	log.Printf("  Metrics: %s", req.Metrics)

	// Process the request (add your business logic here)
	
	// Return a response
	return &pb.ReportResultResponse{
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
	pb.RegisterAggregatorServer(grpcServer, &server{})
	
	log.Printf("gRPC server listening on %s", port)
	
	// Start serving (this blocks)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}