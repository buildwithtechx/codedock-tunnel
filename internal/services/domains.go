package services

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"codedock.run/codedock-tunnel/internal/auth"
	"codedock.run/codedock-tunnel/internal/models"
	"codedock.run/codedock-tunnel/internal/repositories"
)

type DomainService struct {
	domains repositories.DomainRepository
	now     func() time.Time
}

func NewDomainService(domains repositories.DomainRepository) (*DomainService, error) {
	if domains == nil {
		return nil, fmt.Errorf("domain repository is required")
	}
	return &DomainService{domains: domains, now: time.Now}, nil
}

func (s *DomainService) Create(ctx context.Context, organizationID, hostname, method string, tunnelID *string) (string, models.Domain, error) {
	hostname = strings.ToLower(strings.TrimSpace(hostname))
	method = strings.ToLower(strings.TrimSpace(method))
	if organizationID == "" || !validHostname(hostname) || method == "" {
		return "", models.Domain{}, fmt.Errorf("organization, hostname, and verification method are required")
	}
	raw, err := auth.NewToken("cdv", 24)
	if err != nil {
		return "", models.Domain{}, err
	}
	domain := models.Domain{OrganizationID: organizationID, TunnelID: tunnelID, Hostname: hostname, Status: models.DomainStatusPending, VerificationMethod: method, VerificationToken: auth.HashToken(raw), CertificateStatus: "pending"}
	if err := s.domains.Create(ctx, &domain); err != nil {
		return "", models.Domain{}, fmt.Errorf("create domain: %w", err)
	}
	return raw, domain, nil
}

func (s *DomainService) Verify(ctx context.Context, id, token string) error {
	domain, err := s.domains.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("find domain: %w", err)
	}
	if !auth.EqualHash(domain.VerificationToken, token) {
		return fmt.Errorf("invalid domain verification token")
	}
	now := s.now()
	domain.Status = models.DomainStatusVerified
	domain.VerifiedAt = &now
	if err := s.domains.Update(ctx, &domain); err != nil {
		return fmt.Errorf("verify domain: %w", err)
	}
	return nil
}

func validHostname(hostname string) bool {
	if len(hostname) > 253 || !strings.Contains(hostname, ".") || strings.Contains(hostname, "..") || net.ParseIP(hostname) != nil {
		return false
	}
	for _, label := range strings.Split(hostname, ".") {
		if label == "" || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' {
			return false
		}
		for _, char := range label {
			if (char < 'a' || char > 'z') && (char < '0' || char > '9') && char != '-' {
				return false
			}
		}
	}
	return true
}
