package handler

import (
	"time"

	"continuum/api/internal/model"
	"continuum/api/internal/repository"
	"continuum/api/internal/service"

	"github.com/gofiber/fiber/v3"
)

type VaultHandler struct {
	Repo             *repository.Database
	VaultService     *service.VaultService
	LightningService *service.LightningService 
}

func NewVaultHandler(repo *repository.Database, vaultSvc *service.VaultService, lnSvc *service.LightningService) *VaultHandler {
	return &VaultHandler{
		Repo:             repo,
		VaultService:     vaultSvc,
		LightningService: lnSvc,
	}
}

// CreateVault handles POST /api/vaults
func (h *VaultHandler) CreateVault(c fiber.Ctx) error {
	var req model.Vault
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request payload format"})
	}

	// Security Boundary Check
	if req.Alias == "" || req.BeneficiaryPubkey == "" || req.EncryptedPayload == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing critical cryptographic fields"})
	}

	// Enforce standard safe intervals if client forgets to pass one (Default: 30 days)
	if req.CheckInIntervalSeconds <= 0 {
		req.CheckInIntervalSeconds = 2592000 
	}

	// 🛡️ Data Initialization Overrides
	req.Status = "ACTIVE"
	req.LastCheckInAt = time.Now()

	// FIXED: Swapped out c.Context() for c.UserContext() to align with Fiber v3 architecture
	if err := h.Repo.InsertVault(c.UserContext(), &req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":   "SUCCESS",
		"vault_id": req.ID,
		"message":  "Continuum cryptographic vault deployed and activated safely.",
	})
}