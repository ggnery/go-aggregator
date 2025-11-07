package entities

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Artifact family types
type ArtifactFamily string

const (
	ArtifactFamilyModel     ArtifactFamily = "model"
	ArtifactFamilyCode      ArtifactFamily = "code"
	ArtifactFamilyDataset   ArtifactFamily = "dataset"
	ArtifactFamilyTokenizer ArtifactFamily = "tokenizer"
	// Add other families as needed
)

type Artifact struct {
	ArtifactID      uuid.UUID       `json:"artifact_id" db:"artifact_id"`
	Family          ArtifactFamily  `json:"family" db:"family"`
	ResourceLocator string          `json:"resource_locator" db:"resource_locator"`
	IntegrityHash   sql.NullString  `json:"integrity_hash,omitempty" db:"integrity_hash"` // nullable
	RequiredScope   sql.NullString  `json:"required_scope,omitempty" db:"required_scope"` // nullable
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	Annotations     json.RawMessage `json:"annotations" db:"annotations"`
}

type ArtifactVersion struct {
	ArtifactVersionID uuid.UUID       `json:"artifact_version_id" db:"artifact_version_id"`
	ArtifactID        uuid.UUID       `json:"artifact_id" db:"artifact_id"`
	CreatedAt         time.Time       `json:"created_at" db:"created_at"`
	Meta              json.RawMessage `json:"meta" db:"meta"`
}
