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

func (r *TaskRepository) UpdateTaskByTaskResult(taskResult orm.TaskResult) error {
	// Note: attempts_made is NOT incremented here - it's incremented when the lease is created by the Scheduler
	// This function only updates the task status based on the attempt result
	query := `
		UPDATE orchestrator.task
		SET status = CASE 
		        -- If succeeded, mark as completed
		        WHEN $1 = 'succeeded' THEN 'completed'
		        
		        -- If failed/aborted/expired but can retry, set back to pending
		        WHEN $1 IN ('failed', 'aborted', 'expired') 
		             AND attempts_made < max_attempts 
		        THEN 'pending'
		        
		        -- If failed/aborted and max attempts reached, mark as failed
		        WHEN $1 IN ('failed', 'aborted') 
		             AND attempts_made >= max_attempts 
		        THEN 'failed'
		        
		        -- If expired and max attempts reached, mark as timeout
		        WHEN $1 = 'expired' 
		             AND attempts_made >= max_attempts 
		        THEN 'timeout'
		        
		        -- Default: keep current status
		        ELSE status
		    END
		WHERE task_id = $2
	`

	_, err := r.db.Exec(
		query,
		taskResult.Status,
		taskResult.TaskID,
	)
	if err != nil {
		return err
	}

	return nil
}
