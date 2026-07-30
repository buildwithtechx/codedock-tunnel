package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"codedock.run/codedock-tunnel/internal/models"
	"codedock.run/codedock-tunnel/internal/repositories"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$`)

type OrganizationService struct {
	organizations repositories.OrganizationRepository
	billing       *BillingService
}

func (s *OrganizationService) SetBilling(billing *BillingService) { s.billing = billing }

func NewOrganizationService(organizations repositories.OrganizationRepository) (*OrganizationService, error) {
	if organizations == nil {
		return nil, fmt.Errorf("organization repository is required")
	}
	return &OrganizationService{organizations: organizations}, nil
}

func (s *OrganizationService) Create(ctx context.Context, ownerID, name, slug string) (models.Organization, error) {
	name = strings.TrimSpace(name)
	slug = strings.ToLower(strings.TrimSpace(slug))
	if ownerID == "" || name == "" || !slugPattern.MatchString(slug) {
		return models.Organization{}, fmt.Errorf("owner, name, and valid slug are required")
	}
	organization := models.Organization{Name: name, Slug: slug, OwnerID: ownerID, Settings: `{}`}
	if err := s.organizations.Create(ctx, &organization); err != nil {
		return models.Organization{}, fmt.Errorf("create organization: %w", err)
	}
	member := models.OrganizationMember{OrganizationID: organization.ID, UserID: ownerID, Role: models.MemberRoleOwner}
	if err := s.organizations.AddMember(ctx, &member); err != nil {
		return models.Organization{}, fmt.Errorf("add organization owner: %w", err)
	}
	return organization, nil
}

func (s *OrganizationService) AddMember(ctx context.Context, organizationID, userID string, role models.MemberRole) error {
	if organizationID == "" || userID == "" || !validMemberRole(role) {
		return fmt.Errorf("organization, user, and valid role are required")
	}
	if err := s.checkMemberCapacity(ctx, organizationID); err != nil {
		return err
	}
	if err := s.organizations.AddMember(ctx, &models.OrganizationMember{OrganizationID: organizationID, UserID: userID, Role: role}); err != nil {
		return fmt.Errorf("add organization member: %w", err)
	}
	return nil
}

func (s *OrganizationService) checkMemberCapacity(ctx context.Context, organizationID string) error {
	if s.billing != nil {
		plan, _, err := s.billing.Entitlements(ctx, organizationID)
		if err != nil {
			return fmt.Errorf("check member entitlement: %w", err)
		}
		count, err := s.organizations.CountMembers(ctx, organizationID)
		if err != nil {
			return fmt.Errorf("count organization members: %w", err)
		}
		if plan.MaxMembers > 0 && count >= int64(plan.MaxMembers) {
			return fmt.Errorf("organization member limit reached")
		}
	}
	return nil
}

func (s *OrganizationService) Authorize(ctx context.Context, organizationID, userID string, required models.MemberRole) error {
	member, err := s.organizations.FindMember(ctx, organizationID, userID)
	if err != nil {
		return fmt.Errorf("find organization membership: %w", err)
	}
	if memberRoleRank(member.Role) < memberRoleRank(required) {
		return fmt.Errorf("insufficient organization role")
	}
	return nil
}

func validMemberRole(role models.MemberRole) bool {
	return role == models.MemberRoleOwner || role == models.MemberRoleAdmin || role == models.MemberRoleMember || role == models.MemberRoleViewer
}

func memberRoleRank(role models.MemberRole) int {
	switch role {
	case models.MemberRoleOwner:
		return 4
	case models.MemberRoleAdmin:
		return 3
	case models.MemberRoleMember:
		return 2
	case models.MemberRoleViewer:
		return 1
	default:
		return 0
	}
}
