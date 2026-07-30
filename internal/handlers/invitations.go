package handlers

import (
	"fmt"
	"strings"

	"codedock.run/codedock-tunnel/internal/models"
	"codedock.run/codedock-tunnel/internal/services"
	"codedock.run/codedock-tunnel/internal/validation"
	"github.com/gofiber/fiber/v2"
)

type InvitationHandler struct{ invitations *services.InvitationService }

type CreateInvitationRequest struct {
	Email string            `json:"email" validate:"required,email"`
	Role  models.MemberRole `json:"role" validate:"required"`
}

func NewInvitationHandler(invitations *services.InvitationService) (*InvitationHandler, error) {
	if invitations == nil {
		return nil, fmt.Errorf("invitation service is required")
	}
	return &InvitationHandler{invitations: invitations}, nil
}

func (h *InvitationHandler) Create(c *fiber.Ctx) error {
	var input CreateInvitationRequest
	if err := c.BodyParser(&input); err != nil {
		return writeError(c, fiber.StatusBadRequest, fmt.Errorf("decode invitation request: %w", err))
	}
	if err := validation.Struct(input); err != nil {
		return writeError(c, fiber.StatusBadRequest, err)
	}
	userID, err := sessionUserID(c)
	if err != nil {
		return writeError(c, fiber.StatusUnauthorized, err)
	}
	invitation, err := h.invitations.Invite(c.UserContext(), userID, strings.TrimSpace(c.Params("organizationID")), input.Email, input.Role)
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, err)
	}
	return c.Status(fiber.StatusCreated).JSON(invitation)
}

func (h *InvitationHandler) Accept(c *fiber.Ctx) error {
	userID, err := sessionUserID(c)
	if err != nil {
		return writeError(c, fiber.StatusUnauthorized, err)
	}
	if err := h.invitations.Accept(c.UserContext(), userID, c.Query("token")); err != nil {
		return writeError(c, fiber.StatusBadRequest, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
