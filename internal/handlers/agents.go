package handlers

import (
	"fmt"
	"strings"

	"codedock.run/codedock-tunnel/internal/services"
	"github.com/gofiber/fiber/v2"
)

type AgentHandler struct{ agents *services.AgentService }

type RegisterAgentRequest struct {
	Name string `json:"name"`
}

type AgentHeartbeatRequest struct {
	Version  string `json:"version"`
	Hostname string `json:"hostname"`
	Platform string `json:"platform"`
}

func NewAgentHandler(agents *services.AgentService) (*AgentHandler, error) {
	if agents == nil {
		return nil, fmt.Errorf("agent service is required")
	}
	return &AgentHandler{agents: agents}, nil
}

func (h *AgentHandler) Register(c *fiber.Ctx) error {
	var input RegisterAgentRequest
	if err := c.BodyParser(&input); err != nil {
		return writeError(c, fiber.StatusBadRequest, fmt.Errorf("decode agent request: %w", err))
	}
	token, agent, err := h.agents.Register(c.UserContext(), strings.TrimSpace(c.Params("organizationID")), input.Name)
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, err)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"agent": agent, "token": token})
}

func (h *AgentHandler) Heartbeat(c *fiber.Ctx) error {
	var input AgentHeartbeatRequest
	if err := c.BodyParser(&input); err != nil {
		return writeError(c, fiber.StatusBadRequest, fmt.Errorf("decode heartbeat request: %w", err))
	}
	if err := h.agents.Heartbeat(c.UserContext(), c.Params("agentID"), input.Version, input.Hostname, input.Platform); err != nil {
		return writeError(c, fiber.StatusNotFound, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AgentHandler) Revoke(c *fiber.Ctx) error {
	if err := h.agents.Revoke(c.UserContext(), c.Params("agentID")); err != nil {
		return writeError(c, fiber.StatusNotFound, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
