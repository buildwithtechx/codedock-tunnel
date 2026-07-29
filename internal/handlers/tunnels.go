package handlers

import (
	"fmt"
	"strings"

	"codedock.run/codedock-tunnel/internal/models"
	"codedock.run/codedock-tunnel/internal/services"
	"codedock.run/codedock-tunnel/internal/validation"
	"github.com/gofiber/fiber/v2"
)

type TunnelHandler struct{ tunnels *services.TunnelService }

type CreateTunnelRequest struct {
	Name           string                `json:"name" validate:"required,max=120"`
	Protocol       models.TunnelProtocol `json:"protocol" validate:"required"`
	TargetHost     string                `json:"targetHost" validate:"required,max=253"`
	TargetPort     int                   `json:"targetPort" validate:"gte=1,lte=65535"`
	PublicHostname string                `json:"publicHostname,omitempty" validate:"max=253"`
}

type UpdateTunnelStatusRequest struct {
	Status models.TunnelStatus `json:"status"`
}

func NewTunnelHandler(tunnels *services.TunnelService) (*TunnelHandler, error) {
	if tunnels == nil {
		return nil, fmt.Errorf("tunnel service is required")
	}
	return &TunnelHandler{tunnels: tunnels}, nil
}

func (h *TunnelHandler) Create(c *fiber.Ctx) error {
	var input CreateTunnelRequest
	if err := c.BodyParser(&input); err != nil {
		return writeError(c, fiber.StatusBadRequest, fmt.Errorf("decode tunnel request: %w", err))
	}
	if err := validation.Struct(input); err != nil {
		return writeError(c, fiber.StatusBadRequest, err)
	}
	tunnel, err := h.tunnels.Create(c.UserContext(), strings.TrimSpace(c.Params("organizationID")), input.Name, input.Protocol, input.TargetHost, input.TargetPort, input.PublicHostname)
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, err)
	}
	return c.Status(fiber.StatusCreated).JSON(tunnel)
}

func (h *TunnelHandler) SetStatus(c *fiber.Ctx) error {
	var input UpdateTunnelStatusRequest
	if err := c.BodyParser(&input); err != nil {
		return writeError(c, fiber.StatusBadRequest, fmt.Errorf("decode tunnel status request: %w", err))
	}
	if err := h.tunnels.SetStatus(c.UserContext(), c.Params("tunnelID"), input.Status); err != nil {
		return writeError(c, fiber.StatusBadRequest, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *TunnelHandler) Revoke(c *fiber.Ctx) error {
	if err := h.tunnels.Revoke(c.UserContext(), c.Params("tunnelID")); err != nil {
		return writeError(c, fiber.StatusNotFound, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
