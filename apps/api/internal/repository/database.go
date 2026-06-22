package repository

import (
	"context"
	"time"

	"continuum/api/internal/model"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	Db *sqlx.DB
}

// NewDatabase initializes connection configurations and builds a reliable pool
func NewDatabase(dataSourceName string) (*Database, error) {
	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	
	// Hardened pooling configurations optimized for cloud resource targets
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &Database{Db: db}, nil
}

// InsertVault writes the client-side encrypted package to our storage engine
func (d *Database) InsertVault(ctx context.Context, v *model.Vault) error {
	// FIXED: Included updated schema tracking metrics (created_at, updated_at)
	query := `
		INSERT INTO vaults (alias, beneficiary_pubkey, encrypted_payload, check_in_interval_seconds, last_check_in_at, status)
		VALUES (:alias, :beneficiary_pubkey, :encrypted_payload, :check_in_interval_seconds, NOW(), 'ACTIVE')
		RETURNING id, last_check_in_at, status, created_at, updated_at;
	`
	stmt, err := d.Db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Maps returned metadata metrics directly back onto our runtime model reference pointer
	return stmt.GetContext(ctx, v, v)
}

// GetVaultByID finds a specific entry space to let beneficiaries pull or evaluate access locks
func (d *Database) GetVaultByID(ctx context.Context, id string) (*model.Vault, error) {
	var vault model.Vault
	
	// FIXED: Synchronized column lookup parameters to track audit fields completely
	query := `
		SELECT id, alias, beneficiary_pubkey, encrypted_payload, check_in_interval_seconds, last_check_in_at, status, created_at, updated_at 
		FROM vaults 
		WHERE id = $1
	`
	
	err := d.Db.GetContext(ctx, &vault, query, id)
	if err != nil {
		return nil, err
	}
	return &vault, nil
}