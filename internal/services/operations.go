package services

import (
	"context"
	"fmt"
	"time"

	"codedock.run/codedock-tunnel/internal/repositories"
)

type OperationalStore interface {
	ReconcileActiveTunnels(context.Context, string) error
}

type OperationsService struct {
	organizations repositories.OrganizationRepository
	store         OperationalStore
}

func NewOperationsService(organizations repositories.OrganizationRepository, store OperationalStore) (*OperationsService, error) {
	if organizations == nil || store == nil {
		return nil, fmt.Errorf("organization repository and operational store are required")
	}
	return &OperationsService{organizations: organizations, store: store}, nil
}

func (s *OperationsService) Reconcile(ctx context.Context, _ time.Time) error {
	organizations, err := s.organizations.List(ctx)
	if err != nil {
		return fmt.Errorf("list organizations for operational reconciliation: %w", err)
	}
	for _, organization := range organizations {
		if err := s.store.ReconcileActiveTunnels(ctx, organization.ID); err != nil {
			return fmt.Errorf("reconcile organization %s: %w", organization.ID, err)
		}
	}
	return nil
}
