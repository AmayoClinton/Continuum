package service

import (
	"context"
	"fmt"

	"github.com/lightningnetwork/lnd/lnrpc"
)

type LightningService struct {
	lndClient lnrpc.LightningClient
}

func NewLightningService(client lnrpc.LightningClient) *LightningService {
	return &LightningService{lndClient: client}
}

// GenerateProofInvoice creates a customized 1-sat invoice for a specific vault check-in
func (l *LightningService) GenerateProofInvoice(ctx context.Context, vaultID string) (string, error) {
	if l.lndClient == nil {
		// Mock implementation fallback for testing without active local LND nodes
		return fmt.Sprintf("lnbc1u1mockinvoice_for_vault_%s_paid_proves_ownership_and_life_metadata", vaultID), nil
	}

	resp, err := l.lndClient.AddInvoice(ctx, &lnrpc.Invoice{
		Value: 1, // Strict 1 Sat proof cost
		Memo:  "Continuum Proof of Life Check-in | Vault ID: " + vaultID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate LND node invoice: %w", err)
	}

	// FIX: Change PaymentRequest to PaymentReq to align with the lnrpc protobuf definition
	return resp.PaymentReq, nil
}