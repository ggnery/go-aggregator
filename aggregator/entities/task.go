package entities

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	TaskStatusPending     TaskStatus = "pending"
	TaskStatusLeased      TaskStatus = "leased"
	TaskStatusCompleted   TaskStatus = "completed"
	TaskStatusFailed      TaskStatus = "failed"
	TaskStatusTimeout     TaskStatus = "timeout"
	TaskStatusCancelled   TaskStatus = "cancelled"
	TaskStatusQuarantined TaskStatus = "quarantined"
)

type ModelRequirement struct {
	Family    string `json:"family"`
	MinMemory int    `json:"min_mem"`
	Precision string `json:"precision"`
	// Add other fields as needed
}


type Task struct {
	TaskID        uuid.UUID          `json:"task_id" db:"task_id"`
	ParentJob     uuid.UUID          `json:"parent_job" db:"parent_job"`
	Status        TaskStatus         `json:"status" db:"status"`
	CreatedAt     time.Time          `json:"created_at" db:"created_at"`
	ReadyAt       time.Time          `json:"ready_at" db:"ready_at"`
	Priority      int                `json:"priority" db:"priority"`
	AttemptsMade  int                `json:"attempts_made" db:"attempts_made"`
	MaxAttempts   int                `json:"max_attempts" db:"max_attempts"`
	ModelsNeeded  []ModelRequirement `json:"models_needed" db:"models_needed"`
	InputRef      json.RawMessage    `json:"input_ref,omitempty" db:"input_ref"` // nullable
	CodeReference sql.NullString     `json:"code_reference,omitempty" db:"code_reference"` // nullable
	DeadlineAt    sql.NullTime       `json:"deadline_at,omitempty" db:"deadline_at"` // nullable
	Annotations   json.RawMessage    `json:"annotations" db:"annotations"`
}

//=========================================================================================================================

type TaskAttemptStatus string

const (
	TaskAttemptStatusLeased    TaskAttemptStatus = "leased"
	TaskAttemptStatusRunning   TaskAttemptStatus = "running"
	TaskAttemptStatusSucceeded TaskAttemptStatus = "succeeded"
	TaskAttemptStatusFailed    TaskAttemptStatus = "failed"
	TaskAttemptStatusAborted   TaskAttemptStatus = "aborted"
	TaskAttemptStatusExpired   TaskAttemptStatus = "expired"
)

type TaskAttempt struct {
	AttemptID      int64             `json:"attempt_id" db:"attempt_id"`
	TaskID         uuid.UUID         `json:"task_id" db:"task_id"`
	AttemptNo      int               `json:"attempt_no" db:"attempt_no"`
	SessionID      *uuid.UUID        `json:"session_id,omitempty" db:"session_id"` // nullable
	NodeID         *uuid.UUID        `json:"node_id,omitempty" db:"node_id"`        // nullable
	LeaseToken     uuid.UUID         `json:"lease_token" db:"lease_token"`
	LeaseExpiresAt time.Time         `json:"lease_expires_at" db:"lease_expires_at"`
	LeasedAt       time.Time         `json:"leased_at" db:"leased_at"`
	Status         TaskAttemptStatus `json:"status" db:"status"`
	StartedAt      sql.NullTime      `json:"started_at,omitempty" db:"started_at"`       // nullable
	FinishedAt     sql.NullTime      `json:"finished_at,omitempty" db:"finished_at"`     // nullable
	ResultRef      sql.NullString    `json:"result_ref,omitempty" db:"result_ref"`       // nullable
	ErrorCode      sql.NullString    `json:"error_code,omitempty" db:"error_code"`       // nullable
	ErrorMessage   sql.NullString    `json:"error_message,omitempty" db:"error_message"` // nullable
	Metrics        json.RawMessage   `json:"metrics" db:"metrics"`
}