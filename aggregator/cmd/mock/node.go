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
	// Using data from seed file V2__seed_data.sql - simulating completion of the running task
	// Task: 80000000-0000-0000-0000-000000000012 (part of Job 2, batch inference)
	// Attempt: 5 (currently leased and running in seed data)
	request := &aggregatorpb.ReportResultRequest{
		TaskId:     "80000000-0000-0000-0000-000000000012",
		AttemptId:  5,
		LeaseToken: "90000000-0000-0000-0000-000000000012",
		Status:     "succeeded",
		ResultRef:  "s3://results/task-80000000-0000-0000-0000-000000000012.json",
		Metrics:    `{"duration_ms": 450000, "tokens_generated": 95, "batch_size": 10, "memory_mb": 2048}`,
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
