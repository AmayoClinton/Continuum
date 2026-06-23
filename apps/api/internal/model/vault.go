package model

import (
	"time"

	"github.com/lib/pq"
)

type VaultStatus string

const (
	StatusActive  VaultStatus = "ACTIVE"
	StatusDormant VaultStatus = "DORMANT"
)

type Vault struct {
	// omitempty ensures that client creation payloads do not accidentally override 
	// backend database UUID auto-generation sequences with empty string fields.
	ID                     string      `db:"id" json:"id,omitempty"`
	Alias                  string      `db:"alias" json:"alias"`
	BeneficiaryPubkey      string      `db:"beneficiary_pubkey" json:"beneficiary_pubkey"`
	EncryptedPayload       string      `db:"encrypted_payload" json:"encrypted_payload"`
	CheckInIntervalSeconds int         `db:"check_in_interval_seconds" json:"check_in_interval_seconds"`
	LastCheckInAt          time.Time   `db:"last_check_in_at" json:"last_check_in_at"`
	Status                 VaultStatus `db:"status" json:"status,omitempty"`
	
	// AUDIT TELEMETRY: Necessary additions for enterprise ledger maintenance
	CreatedAt              time.Time   `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt              time.Time   `db:"updated_at" json:"updated_at,omitempty"`
}
