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
	
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &Database{Db: db}, nil
}

// InsertVault writes the client-side encrypted package to our storage engine
func (d *Database) InsertVault(ctx context.Context, v *model.Vault) error {
	query := `
		INSERT INTO vaults (alias, beneficiary_pubkey, encrypted_payload, check_in_interval_seconds, last_check_in_at, status)
		VALUES (:alias, :beneficiary_pubkey, :encrypted_payload, :check_in_interval_seconds, NOW(), 'ACTIVE')
		RETURNING id, last_check_in_at, status;
	`
	stmt, err := d.Db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.GetContext(ctx, v, v)
}

// GetVaultByID finds a specific entry space to let beneficiaries pull or evaluate access locks
func (d *Database) GetVaultByID(ctx context.Context, id string) (*model.Vault, error) {
	var vault model.Vault
	query := `SELECT id, alias, beneficiary_pubkey, encrypted_payload, check_in_interval_seconds, last_check_in_at, status FROM vaults WHERE id = $1`
	
	err := d.Db.GetContext(ctx, &vault, query, id)
	if err != nil {
		return nil, err
	}
	return &vault, nil
}