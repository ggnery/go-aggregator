package repository

import (
	"aggregator/repository/orm"
	"database/sql"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) UpdateTaskAttemptbyTaskResult(taskResult orm.TaskResult) error {
	query := `
		UPDATE orchestrator.task_attempt
		SET status = $1,
		    finished_at = now(),
		    result_ref = $2,
		    metrics = $3
		WHERE attempt_id = $4
		  AND task_id = $5
		  AND lease_token = $6
	`

	_, err := r.db.Exec(
		query,
		taskResult.Status,
		taskResult.ResultRef,
		taskResult.Metrics,
		taskResult.AttemptID,
		taskResult.TaskID,
		taskResult.LeaseToken,
	)
	if err != nil {
		return err
	}

	return nil
}
