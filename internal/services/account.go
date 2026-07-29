package services

import (
	"context"
	"fmt"
	"time"

	"codedock.run/codedock-tunnel/internal/repositories"
)

type AccountService struct {
	users         repositories.UserRepository
	organizations repositories.OrganizationRepository
	now           func() time.Time
}

func NewAccountService(users repositories.UserRepository, organizations repositories.OrganizationRepository) (*AccountService, error) {
	if users == nil || organizations == nil {
		return nil, fmt.Errorf("account repositories are required")
	}
	return &AccountService{users: users, organizations: organizations, now: time.Now}, nil
}

func (s *AccountService) Delete(ctx context.Context, userID string) error {
	owned, err := s.organizations.ListOwned(ctx, userID)
	if err != nil {
		return fmt.Errorf("check owned organizations: %w", err)
	}
	if len(owned) > 0 {
		return fmt.Errorf("transfer organization ownership before deleting account")
	}
	if err := s.users.Delete(ctx, userID, s.now()); err != nil {
		return fmt.Errorf("delete account: %w", err)
	}
	return nil
}

func (s *AccountService) TransferOwnership(ctx context.Context, organizationID, currentOwnerID, newOwnerID string) error {
	if organizationID == "" || currentOwnerID == "" || newOwnerID == "" || currentOwnerID == newOwnerID {
		return fmt.Errorf("valid ownership transfer users are required")
	}
	if err := s.organizations.TransferOwnership(ctx, organizationID, currentOwnerID, newOwnerID); err != nil {
		return fmt.Errorf("transfer ownership: %w", err)
	}
	return nil
}
