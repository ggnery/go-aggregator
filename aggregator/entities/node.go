package entities

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Node type categories
type NodeType string

const (
	NodeTypeGPUServer NodeType = "gpu_server"
	NodeTypeEdge      NodeType = "edge"
	NodeTypeLaptop    NodeType = "laptop"
	NodeTypeBrowser   NodeType = "browser"
)

// Node session status
type NodeSessionStatus string

const (
	NodeSessionStatusActive   NodeSessionStatus = "active"
	NodeSessionStatusDraining NodeSessionStatus = "draining"
	NodeSessionStatusLost     NodeSessionStatus = "lost"
	NodeSessionStatusEnded    NodeSessionStatus = "ended"
)

// NodeDescription represents hardware specifications
type NodeDescription struct {
	CPUCores int    `json:"cpu_cores,omitempty"`
	GPUModel string `json:"gpu_model,omitempty"`
	VRAMMB   int64  `json:"vram_mb,omitempty"`
	RAMMB    int64  `json:"ram_mb,omitempty"`
	DiskMB   int64  `json:"disk_mb,omitempty"`
	// Add other fields as needed
}

type Node struct {
	NodeID             uuid.UUID       `json:"node_id" db:"node_id"`
	EnrolledAt         time.Time       `json:"enrolled_at" db:"enrolled_at"`
	AgentVersion       string          `json:"agent_version" db:"agent_version"`
	NodeType           NodeType        `json:"node_type" db:"node_type"`
	OwnerAccountID     *uuid.UUID      `json:"owner_account_id,omitempty" db:"owner_account_id"`         // nullable
	OperatingAccountID *uuid.UUID      `json:"operating_account_id,omitempty" db:"operating_account_id"` // nullable
	Description        NodeDescription `json:"description" db:"description"`
	Region             sql.NullString  `json:"region,omitempty" db:"region"` // nullable
}

type NodeSession struct {
	SessionID       uuid.UUID         `json:"session_id" db:"session_id"`
	NodeID          uuid.UUID         `json:"node_id" db:"node_id"`
	StartedAt       time.Time         `json:"started_at" db:"started_at"`
	EndedAt         sql.NullTime      `json:"ended_at,omitempty" db:"ended_at"` // nullable
	Status          NodeSessionStatus `json:"status" db:"status"`
	Controllable    bool              `json:"controllable" db:"controllable"`
	LastHeartbeatAt sql.NullTime      `json:"last_heartbeat_at,omitempty" db:"last_heartbeat_at"` // nullable
	RTTMS           sql.NullInt64     `json:"rtt_ms,omitempty" db:"rtt_ms"`                       // nullable
	HealthScore     sql.NullFloat64   `json:"health_score,omitempty" db:"health_score"`           // nullable
	Metrics         json.RawMessage   `json:"metrics" db:"metrics"`
}

type NodeCapability struct {
	NodeID uuid.UUID       `json:"node_id" db:"node_id"`
	Cap    json.RawMessage `json:"cap" db:"cap"`
}

type NodeCacheEntry struct {
	NodeID            uuid.UUID     `json:"node_id" db:"node_id"`
	ArtifactVersionID uuid.UUID     `json:"artifact_version_id" db:"artifact_version_id"`
	AcquiredAt        time.Time     `json:"acquired_at" db:"acquired_at"`
	LastUsedAt        sql.NullTime  `json:"last_used_at,omitempty" db:"last_used_at"` // nullable
	SizeBytes         sql.NullInt64 `json:"size_bytes,omitempty" db:"size_bytes"`     // nullable
}
