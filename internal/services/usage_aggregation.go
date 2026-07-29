package services

import (
	"context"
	"fmt"
	"time"

	"codedock.run/codedock-tunnel/internal/repositories"
)

type UsageAggregationService struct {
	organizations repositories.OrganizationRepository
	usage         *UsageService
	now           func() time.Time
}

func NewUsageAggregationService(organizations repositories.OrganizationRepository, usage *UsageService) (*UsageAggregationService, error) {
	if organizations == nil || usage == nil {
		return nil, fmt.Errorf("organization repository and usage service are required")
	}
	return &UsageAggregationService{organizations: organizations, usage: usage, now: time.Now}, nil
}

func (s *UsageAggregationService) Aggregate(ctx context.Context, now time.Time) error {
	periodEnd := now.UTC().Truncate(time.Hour)
	periodStart := periodEnd.Add(-time.Hour)
	organizations, err := s.organizations.List(ctx)
	if err != nil {
		return fmt.Errorf("list organizations for usage aggregation: %w", err)
	}
	for _, organization := range organizations {
		if _, err := s.usage.Aggregate(ctx, organization.ID, periodStart, periodEnd); err != nil {
			return fmt.Errorf("aggregate organization %s: %w", organization.ID, err)
		}
	}
	return nil
}
