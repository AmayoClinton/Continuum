package handler

import (
	"continuum/api/internal/repository"
	"github.com/gofiber/fiber/v3"
)

type ProofHandler struct {
	Repo *repository.Database
}

func NewProofHandler(repo *repository.Database) *ProofHandler {
	return &ProofHandler{Repo: repo}
}

// SimulateTimeWarp handles POST /api/vaults/:id/warp
func (h *ProofHandler) SimulateTimeWarp(c fiber.Ctx) error {
	vaultID := c.Params("id")
	if vaultID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing target vault ID"})
	}

	// Forcibly set the check-in time back 100 days into the past
	query := `
		UPDATE vaults 
		SET last_check_in_at = NOW() - INTERVAL '100 days' 
		WHERE id = $1;
	`
	_, err := h.Repo.Db.ExecContext(c.Context(), query, vaultID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "SUCCESS",
		"message": "Time warped 100 days back. Background loop will flag vault DORMANT on next tick.",
	})
}