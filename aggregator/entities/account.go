package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	AccountID     uuid.UUID      `json:"account_id" db:"account_id"`
	WalletAddress sql.NullString `json:"wallet_address,omitempty" db:"wallet_address"` // nullable
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
}
