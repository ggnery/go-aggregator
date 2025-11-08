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

func (s *ReportResultService) ReportResult(taskResult orm.TaskResult) error {
	//TODO: Implement the report result service

	err := s.taskRepository.UpdateTaskAttemptbyTaskResult(taskResult)
	if err != nil {
		return err
	}

	err = s.taskRepository.UpdateTaskByTaskResult(taskResult)
	if err != nil {
		return err
	}

	err = s.aggregatorRepository.InsertAggregationPart(taskResult)
	if err != nil {
		return err
	}

	return nil
}
