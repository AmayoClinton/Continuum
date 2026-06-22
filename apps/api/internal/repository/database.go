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
	query := `
		INSERT INTO vaults (
			alias, 
			beneficiary_pubkey, 
			encrypted_payload, 
			check_in_interval_seconds, 
			last_check_in_at, 
			status
		) VALUES ($1, $2, $3, $4, NOW(), 'ACTIVE')
		RETURNING id, last_check_in_at, status, created_at, updated_at;
	`
	
	// Explicitly query and map the database-generated metadata metrics 
	// directly onto our runtime model reference pointer fields.
	err := d.Db.GetContext(ctx, v, query, 
		v.Alias, 
		v.BeneficiaryPubkey, 
		v.EncryptedPayload, 
		v.CheckInIntervalSeconds,
	)
	
	return err
}

// GetVaultByID finds a specific entry space to let beneficiaries pull or evaluate access locks
func (d *Database) GetVaultByID(ctx context.Context, id string) (*model.Vault, error) {
	var vault model.Vault
	
	query := `
		SELECT 
			id, 
			alias, 
			beneficiary_pubkey, 
			encrypted_payload, 
			check_in_interval_seconds, 
			last_check_in_at, 
			status, 
			created_at, 
			updated_at 
		FROM vaults 
		WHERE id = $1
	`
	
	err := d.Db.GetContext(ctx, &vault, query, id)
	if err != nil {
		return nil, err
	}
	return &vault, nil
}