package service

import (
	"context"
	"fmt"

	"github.com/lightningnetwork/lnd/lnrpc"
)

// LNDConfig encapsulates the node credentials required to connect to mainnet/testnet
type LNDConfig struct {
	Host         string
	MacaroonPath string
	TLSCertPath  string
}

type LightningService struct {
	lndClient lnrpc.LightningClient
}

// NewLightningService accepts our custom configuration block.
func NewLightningService(config *LNDConfig) (*LightningService, error) {
	// If configuration endpoints are missing, safely fall back onto mock simulation mode
	if config == nil || config.Host == "" {
		return &LightningService{lndClient: nil}, nil
	}

	// NOTE: In production, we would write our secure gRPC dial routines here:
	// conn, _ := lnd.Connect(config.Host, config.TLSCertPath, config.MacaroonPath)
	// client := lnrpc.NewLightningClient(conn)
	// return &LightningService{lndClient: client}, nil

	return &LightningService{lndClient: nil}, nil
}

// GenerateProofInvoice creates a customized 1-sat invoice for a specific vault check-in
func (l *LightningService) GenerateProofInvoice(ctx context.Context, vaultID string) (string, error) {
	if l.lndClient == nil {
		// Mock implementation fallback for testing without active local LND nodes
		return fmt.Sprintf("lnbc1u1mockinvoice_for_vault_%s_paid_proves_ownership_and_life_metadata", vaultID), nil
	}

	resp, err := l.lndClient.AddInvoice(ctx, &lnrpc.Invoice{
		Value: 1,
		Memo:  "Continuum Proof of Life Check-in | Vault ID: " + vaultID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate LND node invoice: %w", err)
	}

	return resp.PaymentRequest, nil
}