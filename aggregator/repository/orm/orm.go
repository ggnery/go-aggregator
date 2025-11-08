package orm

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
)

type TaskResult struct {
	AttemptID int64
	TaskID uuid.UUID
	LeaseToken uuid.UUID
	Status string
	ResultRef sql.NullString
	Metrics json.RawMessage
}