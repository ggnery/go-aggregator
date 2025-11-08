package repository

import (
	"aggregator/repository/orm"
	"database/sql"
	"strconv"

	"github.com/google/uuid"
)

type AggregatorRepository struct {
	db *sql.DB
}

func NewAggregatorRepository(db *sql.DB) *AggregatorRepository {
	return &AggregatorRepository{db: db}
}

func (r *AggregatorRepository) InsertAggregationPartByReportResult(reportResult orm.ReportResult) error {
	// Get the job_id from the task
	var jobID uuid.UUID
	err := r.db.QueryRow(`
		SELECT parent_job FROM orchestrator.task WHERE task_id = $1
	`, reportResult.TaskID).Scan(&jobID)
	if err != nil {
		return err
	}

	// Construct part_key as a stable dedupe key
	// Format: task_id::attempt_id
	// TODO: include partition_index, chunk_id, checksum
	partKey := reportResult.TaskID.String() + "::" + strconv.FormatInt(reportResult.AttemptID, 10)

	// Get payload_ref from ResultRef
	payloadRef := ""
	if reportResult.ResultRef.Valid {
		payloadRef = reportResult.ResultRef.String
	}

	// Determine merge_state based on task attempt status
	// - 'pending': successful result waiting to be merged
	// - 'corrupt': failed/aborted/expired attempts that need re-queue
	// - 'skipped': duplicate delivery (handled by ON CONFLICT)
	// - 'merged': set later by merge service after successful merge
	mergeState := "pending"
	if reportResult.Status == "failed" || reportResult.Status == "aborted" || reportResult.Status == "expired" {
		mergeState = "corrupt"
	}

	// Begin transaction to ensure atomicity
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO orchestrator.aggregation_part (
			job_id, 
			part_key, 
			task_id, 
			attempt_id, 
			payload_ref, 
			merge_state
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (job_id, part_key) DO NOTHING
	`
	// Insert aggregation_part with deduplication
	_, err = tx.Exec(query, jobID, partKey, reportResult.TaskID, reportResult.AttemptID, payloadRef, mergeState)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *AggregatorRepository) UpdateAggregationByReportResult(reportResult orm.ReportResult) error {
    // Only increment if succeeded
    if reportResult.Status != "succeeded" {
        return nil
    }

    var jobID uuid.UUID
    err := r.db.QueryRow(`
        SELECT parent_job FROM orchestrator.task WHERE task_id = $1
    `, reportResult.TaskID).Scan(&jobID)
    if err != nil {
        return err
    }

    _, err = r.db.Exec(`
        UPDATE orchestrator.aggregation
        SET received_parts = received_parts + 1,
            last_part_merged_at = now()
        WHERE job_id = $1
    `, jobID)
    return err
}