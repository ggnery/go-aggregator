package entities

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Aggregation strategy types
type AggregationStrategy string

const (
	AggregationStrategyConcat    AggregationStrategy = "concat"
	AggregationStrategyReduce    AggregationStrategy = "reduce"
	AggregationStrategyMapReduce AggregationStrategy = "map-reduce"
	AggregationStrategyMajority  AggregationStrategy = "majority"
	AggregationStrategyCustom    AggregationStrategy = "custom"
)

// Merge status for aggregation
type MergeStatus string

const (
	MergeStatusIdle      MergeStatus = "idle"
	MergeStatusMerging   MergeStatus = "merging"
	MergeStatusFinalized MergeStatus = "finalized"
	MergeStatusFailed    MergeStatus = "failed"
)

// Merge state for individual parts
type MergeState string

const (
	MergeStatePending  MergeState = "pending"
	MergeStateMerged   MergeState = "merged"
	MergeStateSkipped  MergeState = "skipped"
	MergeStateCorrupt  MergeState = "corrupt"
)

type Aggregation struct {
	JobID            uuid.UUID           `json:"job_id" db:"job_id"`
	ExpectedParts    int                 `json:"expected_parts" db:"expected_parts"`
	ReceivedParts    int                 `json:"received_parts" db:"received_parts"`
	Strategy         AggregationStrategy `json:"strategy" db:"strategy"`
	MergeStatus      MergeStatus         `json:"merge_status" db:"merge_status"`
	FinalRef         sql.NullString      `json:"final_ref,omitempty" db:"final_ref"`                       // nullable
	LastPartMergedAt sql.NullTime        `json:"last_part_merged_at,omitempty" db:"last_part_merged_at"` // nullable
	Annotations      json.RawMessage     `json:"annotations" db:"annotations"`
}

type AggregationPart struct {
	JobID      uuid.UUID  `json:"job_id" db:"job_id"`
	PartKey    string     `json:"part_key" db:"part_key"`
	TaskID     *uuid.UUID `json:"task_id,omitempty" db:"task_id"`       // nullable
	AttemptID  *int64     `json:"attempt_id,omitempty" db:"attempt_id"` // nullable
	ReceivedAt time.Time  `json:"received_at" db:"received_at"`
	PayloadRef string     `json:"payload_ref" db:"payload_ref"`
	MergeState MergeState `json:"merge_state" db:"merge_state"`
}