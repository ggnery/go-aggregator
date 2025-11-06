package controller

import (
	"context" // handle request contexts
	"fmt"
	"log"
	aggregatorpb "aggregator/proto/aggregator/v1"
)

type AggregatorServer struct {
	aggregatorpb.UnimplementedAggregatorServer
}

//Report Result Controller
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
		Message: fmt.Sprintf("Task %s result received successfully", req.TaskId),
	}, nil
}
