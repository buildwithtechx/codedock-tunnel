package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"codedock.run/codedock-tunnel/internal/auth"
	"codedock.run/codedock-tunnel/internal/models"
	"codedock.run/codedock-tunnel/internal/repositories"
)

type AgentService struct {
	agents repositories.AgentRepository
	now    func() time.Time
}

func NewAgentService(agents repositories.AgentRepository) (*AgentService, error) {
	if agents == nil {
		return nil, fmt.Errorf("agent repository is required")
	}
	return &AgentService{agents: agents, now: time.Now}, nil
}

func (s *AgentService) Register(ctx context.Context, organizationID, name string) (string, models.Agent, error) {
	if organizationID == "" || strings.TrimSpace(name) == "" {
		return "", models.Agent{}, fmt.Errorf("organization and agent name are required")
	}
	raw, err := auth.NewToken("cda", 32)
	if err != nil {
		return "", models.Agent{}, err
	}
	agent := models.Agent{OrganizationID: organizationID, Name: strings.TrimSpace(name), TokenHash: auth.HashToken(raw), Status: models.AgentStatusPending, Metadata: `{}`}
	if err := s.agents.Create(ctx, &agent); err != nil {
		return "", models.Agent{}, fmt.Errorf("register agent: %w", err)
	}
	return raw, agent, nil
}

func (s *AgentService) Authenticate(ctx context.Context, raw string) (models.Agent, error) {
	if strings.TrimSpace(raw) == "" {
		return models.Agent{}, fmt.Errorf("agent token is required")
	}
	agent, err := s.agents.FindByTokenHash(ctx, auth.HashToken(raw))
	if err != nil {
		return models.Agent{}, fmt.Errorf("find agent: %w", err)
	}
	if agent.RevokedAt != nil || agent.Status == models.AgentStatusRevoked {
		return models.Agent{}, fmt.Errorf("agent is revoked")
	}
	return agent, nil
}

func (s *AgentService) Heartbeat(ctx context.Context, id, version, hostname, platform string) error {
	if err := s.agents.Touch(ctx, id, s.now(), version, hostname, platform); err != nil {
		return fmt.Errorf("heartbeat agent: %w", err)
	}
	return nil
}

func (s *AgentService) Revoke(ctx context.Context, id string) error {
	if err := s.agents.Revoke(ctx, id, s.now()); err != nil {
		return fmt.Errorf("revoke agent: %w", err)
	}
	return nil
}
