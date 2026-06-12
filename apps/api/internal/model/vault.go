package model

import "time"

type VaultStatus string

const (
	StatusActive  VaultStatus = "ACTIVE"
	StatusDormant VaultStatus = "DORMANT"
)

type Vault struct {
	ID                     string      `db:"id" json:"id"`
	Alias                  string      `db:"alias" json:"alias"`
	BeneficiaryPubkey      string      `db:"beneficiary_pubkey" json:"beneficiary_pubkey"`
	EncryptedPayload       string      `db:"encrypted_payload" json:"encrypted_payload"` 
	CheckInIntervalSeconds int         `db:"check_in_interval_seconds" json:"check_in_interval_seconds"`
	LastCheckInAt          time.Time   `db:"last_check_in_at" json:"last_check_in_at"`
	Status                 VaultStatus `db:"status" json:"status"`
}