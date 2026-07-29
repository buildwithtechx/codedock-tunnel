package handlers

import (
	"fmt"
	"strings"

	"codedock.run/codedock-tunnel/internal/services"
	"github.com/gofiber/fiber/v2"
)

type DomainHandler struct{ domains *services.DomainService }

type CreateDomainRequest struct {
	Hostname           string  `json:"hostname"`
	VerificationMethod string  `json:"verificationMethod"`
	TunnelID           *string `json:"tunnelId"`
}

type VerifyDomainRequest struct {
	Token string `json:"token"`
}

func NewDomainHandler(domains *services.DomainService) (*DomainHandler, error) {
	if domains == nil {
		return nil, fmt.Errorf("domain service is required")
	}
	return &DomainHandler{domains: domains}, nil
}

func (h *DomainHandler) Create(c *fiber.Ctx) error {
	var input CreateDomainRequest
	if err := c.BodyParser(&input); err != nil {
		return writeError(c, fiber.StatusBadRequest, fmt.Errorf("decode domain request: %w", err))
	}
	token, domain, err := h.domains.Create(c.UserContext(), strings.TrimSpace(c.Params("organizationID")), input.Hostname, input.VerificationMethod, input.TunnelID)
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, err)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"domain": domain, "verificationToken": token})
}

func (h *DomainHandler) Verify(c *fiber.Ctx) error {
	var input VerifyDomainRequest
	if err := c.BodyParser(&input); err != nil {
		return writeError(c, fiber.StatusBadRequest, fmt.Errorf("decode domain verification request: %w", err))
	}
	if err := h.domains.Verify(c.UserContext(), c.Params("domainID"), input.Token); err != nil {
		return writeError(c, fiber.StatusUnauthorized, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
