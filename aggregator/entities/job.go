package entities

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// SLA tier levels
type SLATier string

const (
	SLATierLatency    SLATier = "latency"
	SLATierThroughput SLATier = "throughput"
	SLATierBestEffort SLATier = "best_effort"
	SLATierStandard   SLATier = "standard" // default
)

// Job status types
type JobStatus string

const (
	JobStatusQueued      JobStatus = "queued"
	JobStatusRunning     JobStatus = "running"
	JobStatusAggregating JobStatus = "aggregating"
	JobStatusCompleted   JobStatus = "completed"
	JobStatusFailed      JobStatus = "failed"
	JobStatusCancelled   JobStatus = "cancelled"
)

type Job struct {
	JobID          uuid.UUID       `json:"job_id" db:"job_id"`
	OwnerID        uuid.UUID       `json:"owner_id" db:"owner_id"`
	ProjectID      *uuid.UUID      `json:"project_id,omitempty" db:"project_id"`       // nullable
	Priority       int             `json:"priority" db:"priority"`
	SLATier        SLATier         `json:"sla_tier" db:"sla_tier"`
	Status         JobStatus       `json:"status" db:"status"`
	SubmittedAt    time.Time       `json:"submitted_at" db:"submitted_at"`
	StartedAt      sql.NullTime    `json:"started_at,omitempty" db:"started_at"`       // nullable
	FinishedAt     sql.NullTime    `json:"finished_at,omitempty" db:"finished_at"`     // nullable
	OutputRef      sql.NullString  `json:"output_ref,omitempty" db:"output_ref"`       // nullable
	ParentJob      *uuid.UUID      `json:"parent_job,omitempty" db:"parent_job"`       // nullable
	IdempotencyKey sql.NullString  `json:"idempotency_key,omitempty" db:"idempotency_key"` // nullable
	CodeReference  sql.NullString  `json:"code_reference,omitempty" db:"code_reference"`   // nullable
	GasFeeEst      sql.NullString  `json:"gas_fee_est,omitempty" db:"gas_fee_est"`         // nullable numeric(38,18)
	GasFeeActual   sql.NullString  `json:"gas_fee_actual,omitempty" db:"gas_fee_actual"`   // nullable numeric(38,18)
	DeadlineAt     sql.NullTime    `json:"deadline_at,omitempty" db:"deadline_at"`         // nullable
	Annotations    json.RawMessage `json:"annotations" db:"annotations"`
}