package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	aggregatorpb "aggregator/api/proto/aggregator/v1"
)

func main() {
	// Set up a connection to the server
	serverAddr := "localhost:50051"

	// Create connection with insecure credentials (for development)
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Create the Aggregator client
	client := aggregatorpb.NewAggregatorClient(conn)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Prepare the ReportResult request
	request := &aggregatorpb.ReportResultRequest{
		TaskId:     "123e4567-e89b-12d3-a456-426614174000",
		AttemptId:  1,
		LeaseToken: "23e45671-e89b-12d3-a456-426614174001",
		Status:     "succeeded",
		ResultRef:  "s3://results/output-task-123.txt",
		Metrics:    `{"duration_ms": 8234, "tokens": 156, "memory_mb": 512}`,
	}

	log.Println("Calling ReportResult RPC...")
	log.Printf("Request: TaskID=%s, AttemptID=%d, Status=%s",
		request.TaskId, request.AttemptId, request.Status)

	// Call the ReportResult RPC
	response, err := client.ReportResult(ctx, request)
	if err != nil {
		log.Fatalf("ReportResult RPC failed: %v", err)
	}

	// Handle the response
	log.Println("\n=== Response Received ===")
	log.Printf("Acknowledged: %v", response.Acknowledged)
	log.Printf("Message: %s", response.Message)
	log.Println("=== Client completed successfully ===")
}
