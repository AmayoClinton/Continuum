package repository

import (
	"context"
	"database/sql"
	"time"

	"continuum/api/internal/model"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

// ListVaults returns vaults in the order a demo operator expects to scan them.
func (d *Database) ListVaults(ctx context.Context) ([]model.Vault, error) {
	var vaults []model.Vault
	query := `
		SELECT
			id,
			alias,
			beneficiary_pubkey,
			encrypted_payload,
			check_in_interval_seconds,
			last_check_in_at,
			status,
			multisig_required,
			multisig_pubkeys,
			multisig_address,
			multisig_redeem_script,
			multisig_descriptor,
			multisig_network
		FROM vaults
		ORDER BY last_check_in_at DESC;
	`
	if err := d.Db.SelectContext(ctx, &vaults, query); err != nil {
		return nil, err
	}
	return vaults, nil
}

// TouchVault records a successful proof-of-life heartbeat.
func (d *Database) TouchVault(ctx context.Context, id string) error {
	result, err := d.Db.ExecContext(ctx, `
		UPDATE vaults
		SET last_check_in_at = NOW(),
		    status = 'ACTIVE'
		WHERE id = $1;
	`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (d *Database) UpdateVaultInterval(ctx context.Context, id string, seconds int) error {
	result, err := d.Db.ExecContext(ctx, `
		UPDATE vaults
		SET check_in_interval_seconds = $2
		WHERE id = $1;
	`, id, seconds)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (d *Database) UpdateVaultMultisigPolicy(ctx context.Context, v *model.Vault) error {
	result, err := d.Db.ExecContext(ctx, `
		UPDATE vaults
		SET multisig_required = $2,
		    multisig_pubkeys = $3,
		    multisig_address = $4,
		    multisig_redeem_script = $5,
		    multisig_descriptor = $6,
		    multisig_network = $7
		WHERE id = $1;
	`,
		v.ID,
		v.MultisigRequired,
		v.MultisigPubkeys,
		v.MultisigAddress,
		v.MultisigRedeemScript,
		v.MultisigDescriptor,
		v.MultisigNetwork,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (d *Database) NormalizeVaultArrays(v *model.Vault) {
	if v.MultisigPubkeys == nil {
		v.MultisigPubkeys = pq.StringArray{}
	}
}
