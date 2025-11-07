package entities

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Entity types for event tracking
type EntityType string

const (
	EntityTypeJob       EntityType = "job"
	EntityTypeTask      EntityType = "task"
	EntityTypeAttempt   EntityType = "attempt"
	EntityTypeNode      EntityType = "node"
	EntityTypeSession   EntityType = "session"
	EntityTypeScheduler EntityType = "scheduler"
)

// Event types for various lifecycle events
type EventType string

const (
	EventTypeSubmitted EventType = "submitted"
	EventTypeLeased    EventType = "leased"
	EventTypeStarted   EventType = "started"
	EventTypeSucceeded EventType = "succeeded"
	EventTypeFailed    EventType = "failed"
	EventTypeHeartbeat EventType = "heartbeat"
	EventTypeReaped    EventType = "reaped"
	// Add other event types as needed
)

// Outbox topic types
type OutboxTopic string

const (
	OutboxTopicMetrics OutboxTopic = "metrics"
	OutboxTopicLogs    OutboxTopic = "logs"
	OutboxTopicWebhook OutboxTopic = "webhook"
)

// Event represents an append-only event log entry (partitioned by month)
type Event struct {
	EventID    int64           `json:"event_id" db:"event_id"`
	TS         time.Time       `json:"ts" db:"ts"`
	EntityType EntityType      `json:"entity_type" db:"entity_type"`
	EntityID   string          `json:"entity_id" db:"entity_id"`
	EventType  EventType       `json:"event_type" db:"event_type"`
	Data       json.RawMessage `json:"data" db:"data"`
}

// OutboxEvent represents an outbox pattern entry for reliable event delivery
type OutboxEvent struct {
	OutboxID    int64           `json:"outbox_id" db:"outbox_id"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	Topic       OutboxTopic     `json:"topic" db:"topic"`
	Payload     json.RawMessage `json:"payload" db:"payload"`
	DeliveredAt sql.NullTime    `json:"delivered_at,omitempty" db:"delivered_at"` // nullable
}
