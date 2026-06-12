// apps/api/internal/service/vault.go
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
	if vault.Alias == "" || vault.BeneficiaryPubkey == "" || vault.EncryptedPayload == "" {
		return errors.New("missing critical cryptographic parameters")
	}
	if vault.CheckInIntervalSeconds <= 0 {
		return errors.New("check-in interval window must be greater than zero")
	}
	if vault.MultisigRequired <= 0 || len(vault.MultisigPubkeys) == 0 || vault.MultisigDescriptor == "" {
		return errors.New("missing multisig recovery policy")
	}
	return s.repo.InsertVault(ctx, vault)
}

// ProcessCheckIn bumps the last seen timestamp of a vault back to NOW(), proving ownership
func (s *VaultService) ProcessCheckIn(ctx context.Context, vaultID string) error {
	return s.repo.TouchVault(ctx, vaultID)
}
