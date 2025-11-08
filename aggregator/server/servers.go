package server

import (
	"aggregator/entities"
	"aggregator/services"

	"context"
	"database/sql"
	"encoding/json"

	aggregatorpb "aggregator/api/proto/aggregator/v1"
	
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func RegisterGRPCServers(grpcServer *grpc.Server) {
	// Register our Aggregator service with the gRPC server
	aggregatorpb.RegisterAggregatorServer(grpcServer, &AggregatorServer{})

	//Register more services here ->
}

type AggregatorServer struct {
	aggregatorpb.UnimplementedAggregatorServer
}

//Report Result server
func (s *AggregatorServer) ReportResult(ctx context.Context, req *aggregatorpb.ReportResultRequest) (*aggregatorpb.ReportResultResponse, error) {
	taskID, err := uuid.Parse(req.TaskId)
	if err != nil {
		return nil, err
	}

	leaseToken, err := uuid.Parse(req.LeaseToken)
	if err != nil {
		return nil, err
	}
	
	taskAttempt := entities.TaskAttempt{
		AttemptID: req.AttemptId,
		TaskID: taskID,
		LeaseToken: leaseToken,
		Status: entities.TaskAttemptStatus(req.Status),
		ResultRef: sql.NullString{String: req.ResultRef},
		Metrics: json.RawMessage(req.Metrics),
	}

	err = services.ReportResultService(taskAttempt)
	if err != nil {
		return nil, err
	}

	return &aggregatorpb.ReportResultResponse{Acknowledged: true, Message: "200"}, nil
}


