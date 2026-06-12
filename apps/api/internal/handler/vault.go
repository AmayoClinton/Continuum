package handler

import (
	"encoding/json"
	"continuum/api/internal/model"
	"continuum/api/internal/repository"
	"continuum/api/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/lib/pq"
)

type VaultHandler struct {
	Repo      *repository.Database
	Service   *service.VaultService
	Lightning *service.LightningService
	Multisig  *service.MultisigService
}

type updateTimerRequest struct {
	CheckInIntervalSeconds int `json:"check_in_interval_seconds"`
}

type addBeneficiaryRequest struct {
	Pubkey string `json:"pubkey"`
}

func NewVaultHandler(repo *repository.Database, vaultService *service.VaultService, lightningService *service.LightningService, multisigService *service.MultisigService) *VaultHandler {
	return &VaultHandler{Repo: repo, Service: vaultService, Lightning: lightningService, Multisig: multisigService}
}

func (h *VaultHandler) ListVaults(c fiber.Ctx) error {
	vaults, err := h.Repo.ListVaults(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "SUCCESS",
		"vaults": vaults,
	})
}

// CreateVault handles POST /api/vaults
func (h *VaultHandler) CreateVault(c fiber.Ctx) error {
	var req model.Vault
	body := c.Body()
	if len(body) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Empty request body"})
	}
	if err := json.Unmarshal(body, &req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request payload format"})
	}

	if req.Alias == "" || req.BeneficiaryPubkey == "" || req.EncryptedPayload == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing critical cryptographic fields"})
	}
	if req.CheckInIntervalSeconds <= 0 {
		req.CheckInIntervalSeconds = 60
	}

	pubkeys := []string(req.MultisigPubkeys)
	if len(pubkeys) == 0 {
		pubkeys = defaultDemoPubkeys(req.BeneficiaryPubkey)
	}
	policy, err := h.Multisig.BuildPolicy(c.Context(), req.MultisigRequired, pubkeys)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	req.MultisigRequired = policy.Required
	req.MultisigPubkeys = pq.StringArray(policy.Pubkeys)
	req.MultisigAddress = policy.Address
	req.MultisigRedeemScript = policy.RedeemScript
	req.MultisigDescriptor = policy.Descriptor
	req.MultisigNetwork = policy.Network

	if err := h.Service.CreateNewVault(c.Context(), &req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"status":   "SUCCESS",
		"vault_id": req.ID,
		"message":  "Continuum cryptographic vault deployed.",
		"multisig": fiber.Map{
			"required":      req.MultisigRequired,
			"pubkeys":       req.MultisigPubkeys,
			"address":       req.MultisigAddress,
			"redeem_script": req.MultisigRedeemScript,
			"descriptor":    req.MultisigDescriptor,
			"network":       req.MultisigNetwork,
		},
	})
}

func (h *VaultHandler) CheckIn(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing target vault ID"})
	}

	if err := h.Service.ProcessCheckIn(c.Context(), id); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "SUCCESS",
		"message": "Proof-of-life heartbeat accepted. Vault is ACTIVE.",
	})
}

func (h *VaultHandler) UpdateTimer(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing target vault ID"})
	}

	var req updateTimerRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid timer payload"})
	}
	if req.CheckInIntervalSeconds < 30 {
		return c.Status(400).JSON(fiber.Map{"error": "Timer must be at least 30 seconds"})
	}

	if err := h.Repo.UpdateVaultInterval(c.Context(), id, req.CheckInIntervalSeconds); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "SUCCESS",
		"message": "Inactivity timer updated.",
	})
}

func (h *VaultHandler) AddBeneficiary(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing target vault ID"})
	}

	var req addBeneficiaryRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid beneficiary payload"})
	}

	vault, err := h.Repo.GetVaultByID(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Target vault space not found"})
	}

	pubkeys := append([]string{}, vault.MultisigPubkeys...)
	pubkeys = append(pubkeys, req.Pubkey)
	policy, err := h.Multisig.BuildPolicy(c.Context(), vault.MultisigRequired, pubkeys)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	vault.MultisigRequired = policy.Required
	vault.MultisigPubkeys = pq.StringArray(policy.Pubkeys)
	vault.MultisigAddress = policy.Address
	vault.MultisigRedeemScript = policy.RedeemScript
	vault.MultisigDescriptor = policy.Descriptor
	vault.MultisigNetwork = policy.Network

	if err := h.Repo.UpdateVaultMultisigPolicy(c.Context(), vault); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "SUCCESS",
		"message": "Beneficiary signer added to multisig recovery policy.",
		"multisig": fiber.Map{
			"required":      vault.MultisigRequired,
			"pubkeys":       vault.MultisigPubkeys,
			"address":       vault.MultisigAddress,
			"redeem_script": vault.MultisigRedeemScript,
			"descriptor":    vault.MultisigDescriptor,
			"network":       vault.MultisigNetwork,
		},
	})
}

func (h *VaultHandler) CreateProofInvoice(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing target vault ID"})
	}

	invoice, err := h.Lightning.GenerateProofInvoice(c.Context(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "SUCCESS",
		"invoice": invoice,
	})
}

func defaultDemoPubkeys(beneficiaryPubkey string) []string {
	return []string{
		beneficiaryPubkey,
		"02aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"03bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
	}
}
