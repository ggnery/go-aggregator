package services

import (
	"aggregator/repository"
	"aggregator/repository/orm"
)

type ReportResultService struct {
	taskRepository *repository.TaskRepository
}

func NewReportResultService(taskRepository *repository.TaskRepository) *ReportResultService {
	return &ReportResultService{taskRepository: taskRepository}
}

func (s *ReportResultService) ReportResult(taskResult orm.TaskResult) error {
	//TODO: Implement the report result service

	err := s.taskRepository.UpdateTaskAttemptbyTaskResult(taskResult)
	if err != nil {
		return err
	}

	return nil
}
