package server

import (
	"aggregator/repository"
	"aggregator/repository/orm"
	"aggregator/services"

	"context"
	"database/sql"
	"encoding/json"

	aggregatorpb "aggregator/api/proto/aggregator/v1"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func RegisterGRPCServers(grpcServer *grpc.Server, db *sql.DB) {
	// Register repositories here
	taskRepository := repository.NewTaskRepository(db)
	aggregatorRepository := repository.NewAggregatorRepository(db)

	// Register services here
	reportResultService := services.NewReportResultService(taskRepository, aggregatorRepository)

	// Register gRPC servers here
	aggregatorpb.RegisterAggregatorServer(grpcServer, &AggregatorServer{ReportResultService: reportResultService})

}

type AggregatorServer struct {
	aggregatorpb.UnimplementedAggregatorServer
	ReportResultService *services.ReportResultService
}

// Report Result gRPC server
func (s *AggregatorServer) ReportResult(ctx context.Context, req *aggregatorpb.ReportResultRequest) (*aggregatorpb.ReportResultResponse, error) {
	taskID, err := uuid.Parse(req.TaskId)
	if err != nil {
		return nil, err
	}

	leaseToken, err := uuid.Parse(req.LeaseToken)
	if err != nil {
		return nil, err
	}

	reportResult := orm.ReportResult{
		AttemptID:  req.AttemptId,
		TaskID:     taskID,
		LeaseToken: leaseToken,
		Status:     req.Status,
		ResultRef:  sql.NullString{String: req.ResultRef, Valid: true},
		Metrics:    json.RawMessage(req.Metrics),
	}

	err = s.ReportResultService.ReportResult(reportResult)
	if err != nil {
		return nil, err
	}

	return &aggregatorpb.ReportResultResponse{Acknowledged: true, Message: "200"}, nil
}
