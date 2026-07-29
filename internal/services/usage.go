package services

import (
	"context"
	"fmt"
	"time"

	"codedock.run/codedock-tunnel/internal/models"
	"codedock.run/codedock-tunnel/internal/repositories"
)

type UsageService struct {
	usage repositories.UsageRepository
}

func NewUsageService(usage repositories.UsageRepository) (*UsageService, error) {
	if usage == nil {
		return nil, fmt.Errorf("usage repository is required")
	}
	return &UsageService{usage: usage}, nil
}

func (s *UsageService) Record(ctx context.Context, event *models.UsageEvent) error {
	if event == nil || event.OrganizationID == "" || event.EventType == "" || event.Bytes < 0 || event.Connections < 0 {
		return fmt.Errorf("invalid usage event")
	}
	if event.OccurredAt.IsZero() {
		event.OccurredAt = time.Now()
	}
	if err := s.usage.CreateEvent(ctx, event); err != nil {
		return fmt.Errorf("record usage event: %w", err)
	}
	return nil
}

func (s *UsageService) Snapshot(ctx context.Context, snapshot *models.UsageSnapshot) error {
	if snapshot == nil || snapshot.OrganizationID == "" || snapshot.PeriodStart.IsZero() || snapshot.PeriodEnd.IsZero() || snapshot.PeriodEnd.Before(snapshot.PeriodStart) {
		return fmt.Errorf("invalid usage snapshot")
	}
	if snapshot.TunnelCount < 0 || snapshot.ActiveConnections < 0 || snapshot.BandwidthBytes < 0 || snapshot.RequestCount < 0 || snapshot.ErrorCount < 0 {
		return fmt.Errorf("usage counters cannot be negative")
	}
	if err := s.usage.UpsertSnapshot(ctx, snapshot); err != nil {
		return fmt.Errorf("save usage snapshot: %w", err)
	}
	return nil
}
