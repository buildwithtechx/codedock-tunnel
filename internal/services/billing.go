package services

import (
	"context"
	"fmt"
	"time"

	"codedock.run/codedock-tunnel/internal/models"
	"codedock.run/codedock-tunnel/internal/repositories"
)

type BillingService struct {
	billing repositories.BillingRepository
	gateway BillingGateway
	now     func() time.Time
}

func NewBillingService(billing repositories.BillingRepository) (*BillingService, error) {
	if billing == nil {
		return nil, fmt.Errorf("billing repository is required")
	}
	return &BillingService{billing: billing, now: time.Now}, nil
}

func (s *BillingService) Entitlements(ctx context.Context, organizationID string) (models.Plan, models.Subscription, error) {
	subscription, err := s.billing.FindSubscription(ctx, organizationID)
	if err != nil {
		return models.Plan{}, models.Subscription{}, fmt.Errorf("find subscription: %w", err)
	}
	if subscription.Status != models.SubscriptionStatusActive && subscription.Status != models.SubscriptionStatusTrialing {
		return models.Plan{}, models.Subscription{}, fmt.Errorf("subscription is not entitled")
	}
	plan, err := s.billing.FindPlan(ctx, subscription.PlanID)
	if err != nil {
		return models.Plan{}, models.Subscription{}, fmt.Errorf("find subscription plan: %w", err)
	}
	if !plan.Active {
		return models.Plan{}, models.Subscription{}, fmt.Errorf("subscription plan is inactive")
	}
	return plan, subscription, nil
}

func (s *BillingService) RecordEvent(ctx context.Context, event *models.BillingEvent) (bool, error) {
	if event == nil || event.Provider == "" || event.ProviderEventID == "" || event.PayloadHash == "" {
		return false, fmt.Errorf("complete billing event is required")
	}
	_, err := s.billing.FindBillingEvent(ctx, event.Provider, event.ProviderEventID)
	if err == nil {
		return false, nil
	}
	if err != repositories.ErrNotFound {
		return false, fmt.Errorf("check billing event: %w", err)
	}
	if err := s.billing.CreateBillingEvent(ctx, event); err != nil {
		return false, fmt.Errorf("record billing event: %w", err)
	}
	return true, nil
}

func (s *BillingService) MarkProcessed(ctx context.Context, eventID string) error {
	if err := s.billing.MarkBillingEventProcessed(ctx, eventID, s.now()); err != nil {
		return fmt.Errorf("mark billing event processed: %w", err)
	}
	return nil
}
