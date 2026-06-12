package handler

import (
	"continuum/api/internal/model"
	"continuum/api/internal/repository"

	"github.com/gofiber/fiber/v3"
)

type VaultHandler struct {
	Repo *repository.Database
}

func NewVaultHandler(repo *repository.Database) *VaultHandler {
	return &VaultHandler{Repo: repo}
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

	if err := h.Repo.InsertVault(c.Context(), &req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":   "SUCCESS",
		"vault_id": req.ID,
		"message":  "Continuum cryptographic vault deployed.",
	})
}