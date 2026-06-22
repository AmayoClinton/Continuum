package service

import (
	"context"
	"errors"

	"continuum/api/internal/model"
	"continuum/api/internal/repository"
)

type VaultService struct {
	repo *repository.Database
}

func NewVaultService(repo *repository.Database) *VaultService {
	return &VaultService{repo: repo}
}

// CreateNewVault orchestrates the database insertion for Alice's encrypted capsule
func (s *VaultService) CreateNewVault(ctx context.Context, vault *model.Vault) error {
	// Basic business layer validation fallback guard
	if vault.CheckInIntervalSeconds <= 0 {
		return errors.New("check-in interval window must be greater than zero")
	}
	return s.repo.InsertVault(ctx, vault)
}

// ProcessCheckIn bumps the last seen timestamp of a vault back to NOW(), proving ownership
func (s *VaultService) ProcessCheckIn(ctx context.Context, vaultID string) error {
	// FIXED: Updated both last_check_in_at and the standard updated_at audit metric
	query := `
		UPDATE vaults 
		SET last_check_in_at = NOW(), 
		    status = 'ACTIVE',
		    updated_at = NOW()
		WHERE id = $1;
	`
	result, err := s.repo.Db.ExecContext(ctx, query, vaultID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("target continuum vault space not found")
	}

	return nil
}