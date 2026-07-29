package handlers

import (
	"fmt"
	"strings"

	"codedock.run/codedock-tunnel/internal/models"
	"codedock.run/codedock-tunnel/internal/services"
	"github.com/gofiber/fiber/v2"
)

type OrganizationHandler struct{ organizations *services.OrganizationService }

type CreateOrganizationRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type AddMemberRequest struct {
	UserID string            `json:"userId"`
	Role   models.MemberRole `json:"role"`
}

func NewOrganizationHandler(organizations *services.OrganizationService) (*OrganizationHandler, error) {
	if organizations == nil {
		return nil, fmt.Errorf("organization service is required")
	}
	return &OrganizationHandler{organizations: organizations}, nil
}

func (h *OrganizationHandler) Create(c *fiber.Ctx) error {
	var input CreateOrganizationRequest
	if err := c.BodyParser(&input); err != nil {
		return writeError(c, fiber.StatusBadRequest, fmt.Errorf("decode organization request: %w", err))
	}
	userID, err := sessionUserID(c)
	if err != nil {
		return writeError(c, fiber.StatusUnauthorized, err)
	}
	organization, err := h.organizations.Create(c.UserContext(), userID, input.Name, input.Slug)
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, err)
	}
	return c.Status(fiber.StatusCreated).JSON(organization)
}

func (h *OrganizationHandler) AddMember(c *fiber.Ctx) error {
	var input AddMemberRequest
	if err := c.BodyParser(&input); err != nil {
		return writeError(c, fiber.StatusBadRequest, fmt.Errorf("decode member request: %w", err))
	}
	organizationID := strings.TrimSpace(c.Params("organizationID"))
	if err := h.organizations.AddMember(c.UserContext(), organizationID, input.UserID, input.Role); err != nil {
		return writeError(c, fiber.StatusBadRequest, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func sessionUserID(c *fiber.Ctx) (string, error) {
	session, ok := c.Locals("session").(models.Session)
	if !ok || session.UserID == "" {
		return "", fmt.Errorf("authenticated session is required")
	}
	return session.UserID, nil
}
