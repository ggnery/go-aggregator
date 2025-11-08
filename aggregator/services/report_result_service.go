package services

import (
	"aggregator/repository"
	"aggregator/repository/orm"
)

type ReportResultService struct {
	taskRepository       *repository.TaskRepository
	aggregatorRepository *repository.AggregatorRepository
}

func NewReportResultService(taskRepository *repository.TaskRepository, aggregatorRepository *repository.AggregatorRepository) *ReportResultService {
	return &ReportResultService{taskRepository: taskRepository, aggregatorRepository: aggregatorRepository}
}

func (s *ReportResultService) ReportResult(reportResult orm.ReportResult) error {
	err := s.aggregatorRepository.InsertAggregationPartByReportResult(reportResult)
	if err != nil {
		return err
	}

	err = s.aggregatorRepository.UpdateAggregationByReportResult(reportResult)
	if err != nil {
		return err
	}	
	
	err = s.taskRepository.UpdateTaskAttemptbyReportResult(reportResult)
	if err != nil {
		return err
	}

	err = s.taskRepository.UpdateTaskByReportResult(reportResult)
	if err != nil {
		return err
	}

	return nil
}
