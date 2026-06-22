package handler

import (
	"time"

	"continuum/api/internal/model"
	"continuum/api/internal/repository"
	"continuum/api/internal/service"

	"github.com/gofiber/fiber/v3"
)

type VaultHandler struct {
	Repo         *repository.Database
	VaultService *service.VaultService
}

func NewVaultHandler(repo *repository.Database, vaultSvc *service.VaultService) *VaultHandler {
	return &VaultHandler{
		Repo:         repo,
		VaultService: vaultSvc,
	}
}

// CreateVault handles POST /api/vaults
func (h *VaultHandler) CreateVault(c fiber.Ctx) error {
	var req model.Vault
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request payload format"})
	}

	if req.Alias == "" || req.BeneficiaryPubkey == "" || req.EncryptedPayload == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing critical cryptographic fields"})
	}

	if req.CheckInIntervalSeconds <= 0 {
		req.CheckInIntervalSeconds = 2592000 // Default to 30 days
	}

	req.Status = "ACTIVE"
	req.LastCheckInAt = time.Now()

	// ✅ FIXED FOR FIBER v3: Use c.Context()
	if err := h.VaultService.CreateNewVault(c.Context(), &req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":   "SUCCESS",
		"vault_id": req.ID,
		"message":  "Continuum cryptographic vault deployed and activated safely.",
	})
}

// RequestCheckInToken handles POST /api/vaults/:id/invoice
func (h *VaultHandler) RequestCheckInToken(c fiber.Ctx) error {
	vaultID := c.Params("id")
	if vaultID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing target vault ID"})
	}

	// ✅ FIXED FOR FIBER v3: Use c.Context()
	invoice, err := h.VaultService.RequestCheckInInvoice(c.Context(), vaultID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":          "PENDING_PAYMENT",
		"payment_request": invoice,
		"value_sats":      1,
		"message":         "Settle this 1-sat invoice via your Lightning wallet to prove life status.",
	})
}

// ConfirmCheckIn handles POST /api/vaults/:id/checkin
func (h *VaultHandler) ConfirmCheckIn(c fiber.Ctx) error {
	vaultID := c.Params("id")
	if vaultID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing target vault ID"})
	}

	var payload struct {
		Preimage string `json:"preimage"`
	}
	if err := c.Bind().Body(&payload); err != nil || payload.Preimage == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing valid payment proof preimage string"})
	}

	// ✅ FIXED FOR FIBER v3: Use c.Context()
	err := h.VaultService.VerifyAndProcessCheckIn(c.Context(), vaultID, payload.Preimage)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "SUCCESS",
		"message": "Proof of life accepted. Vault dead-man switch timer has been reset to NOW().",
	})
}