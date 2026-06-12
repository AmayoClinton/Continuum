package handler

import (
	"continuum/api/internal/repository"
	"github.com/gofiber/fiber/v3"
)

type RecoveryHandler struct {
	Repo *repository.Database
}

func NewRecoveryHandler(repo *repository.Database) *RecoveryHandler {
	return &RecoveryHandler{Repo: repo}
}

// GetVaultStatus handles GET /api/vaults/:id
func (h *RecoveryHandler) GetVaultStatus(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing target vault ID"})
	}

	vault, err := h.Repo.GetVaultByID(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Target vault space not found"})
	}

	// SECURITY CHECK: If owner is active, shield the payload!
	if vault.Status == "ACTIVE" {
		return c.JSON(fiber.Map{
			"id":                  vault.ID,
			"alias":               vault.Alias,
			"status":              vault.Status,
			"last_seen":           vault.LastCheckInAt,
			"payload_locked":      true,
			"message":             "Vault remains cryptographically shielded. Owner verified active.",
			"multisig_required":   vault.MultisigRequired,
			"multisig_pubkeys":    vault.MultisigPubkeys,
			"multisig_address":    vault.MultisigAddress,
			"multisig_descriptor": vault.MultisigDescriptor,
			"multisig_network":    vault.MultisigNetwork,
		})
	}

	// RELEASE CONDITION MET: Owner missing, hand off the ciphertext package
	return c.JSON(fiber.Map{
		"id":                     vault.ID,
		"alias":                  vault.Alias,
		"status":                 vault.Status,
		"last_seen":              vault.LastCheckInAt,
		"payload_locked":         false,
		"encrypted_payload":      vault.EncryptedPayload, // Released for client-side browser decryption
		"multisig_required":      vault.MultisigRequired,
		"multisig_pubkeys":       vault.MultisigPubkeys,
		"multisig_address":       vault.MultisigAddress,
		"multisig_redeem_script": vault.MultisigRedeemScript,
		"multisig_descriptor":    vault.MultisigDescriptor,
		"multisig_network":       vault.MultisigNetwork,
	})
}
