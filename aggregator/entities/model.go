package entities

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Common framework types
type Framework string

const (
	FrameworkPyTorch Framework = "pytorch"
	FrameworkONNX    Framework = "onnx"
	// Add other frameworks as needed
)

// Common precision types
type Precision string

const (
	PrecisionFP16 Precision = "fp16"
	PrecisionFP8  Precision = "fp8"
	PrecisionInt8 Precision = "int8"
	// Add other precisions as needed
)

// Device class types
type DeviceClass string

const (
	DeviceClassGPU DeviceClass = "gpu"
	DeviceClassCPU DeviceClass = "cpu"
)

type Model struct {
	ModelID            uuid.UUID       `json:"model_id" db:"model_id"`
	Family             string          `json:"family" db:"family"`
	Framework          Framework       `json:"framework" db:"framework"`
	CreatedAt          time.Time       `json:"created_at" db:"created_at"`
	RequireAccelerator bool            `json:"require_accelerator" db:"require_accelerator"`
	Meta               json.RawMessage `json:"meta" db:"meta"`
}

type ModelVariant struct {
	VariantID         uuid.UUID     `json:"variant_id" db:"variant_id"`
	ModelID           uuid.UUID     `json:"model_id" db:"model_id"`
	CreatedAt         time.Time     `json:"created_at" db:"created_at"`
	Precision         *Precision    `json:"precision,omitempty" db:"precision"`
	DeviceClass       *DeviceClass  `json:"device_class,omitempty" db:"device_class"`
	MinVRAMMB         sql.NullInt64 `json:"min_vram_mb,omitempty" db:"min_vram_mb"`
	ArtifactVersionID *uuid.UUID    `json:"artifact_version_id,omitempty" db:"artifact_version_id"`
}
