package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"codedock.run/codedock-tunnel/internal/models"
	"codedock.run/codedock-tunnel/internal/repositories"
)

type TunnelService struct {
	tunnels   repositories.TunnelRepository
	billing   *BillingService
	allocator *HostnameAllocator
	now       func() time.Time
}

func NewTunnelService(tunnels repositories.TunnelRepository) (*TunnelService, error) {
	if tunnels == nil {
		return nil, fmt.Errorf("tunnel repository is required")
	}
	return &TunnelService{tunnels: tunnels, now: time.Now}, nil
}

func (s *TunnelService) SetBilling(billing *BillingService) {
	s.billing = billing
}

func (s *TunnelService) SetHostnameAllocator(allocator *HostnameAllocator) {
	s.allocator = allocator
}

func (s *TunnelService) Create(ctx context.Context, organizationID, name string, protocol models.TunnelProtocol, targetHost string, targetPort int, publicHostname string) (models.Tunnel, error) {
	if organizationID == "" || strings.TrimSpace(name) == "" || strings.TrimSpace(targetHost) == "" || !validTunnelProtocol(protocol) || targetPort < 1 || targetPort > 65535 {
		return models.Tunnel{}, fmt.Errorf("invalid tunnel configuration")
	}
	if s.allocator != nil {
		allocated, err := s.allocator.Allocate(ctx, strings.TrimSpace(publicHostname))
		if err != nil {
			return models.Tunnel{}, fmt.Errorf("allocate public hostname: %w", err)
		}
		publicHostname = allocated
	}
	if strings.TrimSpace(publicHostname) == "" {
		return models.Tunnel{}, fmt.Errorf("public hostname is required")
	}
	tunnel := models.Tunnel{OrganizationID: organizationID, Name: strings.TrimSpace(name), Protocol: protocol, Status: models.TunnelStatusCreated, TargetHost: strings.TrimSpace(targetHost), TargetPort: targetPort, PublicHostname: strings.ToLower(strings.TrimSpace(publicHostname)), AccessPolicy: `{}`}
	if s.billing != nil {
		plan, _, err := s.billing.Entitlements(ctx, organizationID)
		if err != nil {
			return models.Tunnel{}, fmt.Errorf("check tunnel entitlement: %w", err)
		}
		count, err := s.tunnels.CountByOrganization(ctx, organizationID)
		if err != nil {
			return models.Tunnel{}, fmt.Errorf("count organization tunnels: %w", err)
		}
		if plan.MaxTunnels > 0 && count >= int64(plan.MaxTunnels) {
			return models.Tunnel{}, fmt.Errorf("tunnel plan limit reached")
		}
	}
	if err := s.tunnels.Create(ctx, &tunnel); err != nil {
		return models.Tunnel{}, fmt.Errorf("create tunnel: %w", err)
	}
	return tunnel, nil
}

func (s *TunnelService) SetStatus(ctx context.Context, id string, status models.TunnelStatus) error {
	if !validTunnelStatus(status) {
		return fmt.Errorf("invalid tunnel status")
	}
	if err := s.tunnels.UpdateStatus(ctx, id, status); err != nil {
		return fmt.Errorf("set tunnel status: %w", err)
	}
	return nil
}

func (s *TunnelService) Touch(ctx context.Context, id string) error {
	if err := s.tunnels.Touch(ctx, id, s.now()); err != nil {
		return fmt.Errorf("touch tunnel: %w", err)
	}
	return nil
}

func (s *TunnelService) Revoke(ctx context.Context, id string) error {
	if err := s.tunnels.Revoke(ctx, id, s.now()); err != nil {
		return fmt.Errorf("revoke tunnel: %w", err)
	}
	return nil
}

func validTunnelProtocol(protocol models.TunnelProtocol) bool {
	return protocol == models.TunnelProtocolHTTP || protocol == models.TunnelProtocolHTTPS || protocol == models.TunnelProtocolTCP || protocol == models.TunnelProtocolUDP
}

func validTunnelStatus(status models.TunnelStatus) bool {
	return status == models.TunnelStatusCreated || status == models.TunnelStatusConnecting || status == models.TunnelStatusActive || status == models.TunnelStatusDisconnected || status == models.TunnelStatusExpired || status == models.TunnelStatusRevoked
}
