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
	secrets SecretProtector
}

type BillingTransition struct {
	Provider              models.BillingProvider
	ProviderSubscription  string
	ProviderCustomer      string
	ProviderProduct       string
	Status                models.SubscriptionStatus
	CurrentPeriodEnd      *time.Time
	CancelAtPeriodEnd     bool
	ProviderAuthorization string
}

func NewBillingService(billing repositories.BillingRepository) (*BillingService, error) {
	if billing == nil {
		return nil, fmt.Errorf("billing repository is required")
	}
	return &BillingService{billing: billing, now: time.Now}, nil
}

func (s *BillingService) SetSecretProtector(protector SecretProtector) { s.secrets = protector }

func (s *BillingService) ApplyTransition(ctx context.Context, transition BillingTransition) error {
	if transition.Provider == "" || transition.ProviderSubscription == "" || transition.Status == "" {
		return fmt.Errorf("complete billing transition is required")
	}
	subscription, err := s.billing.FindSubscriptionByProvider(ctx, transition.Provider, transition.ProviderSubscription)
	if err != nil {
		return fmt.Errorf("find subscription transition: %w", err)
	}
	subscription.Status = transition.Status
	subscription.ProviderCustomerID = transition.ProviderCustomer
	subscription.ProviderProductID = transition.ProviderProduct
	subscription.CurrentPeriodEnd = transition.CurrentPeriodEnd
	subscription.CancelAtPeriodEnd = transition.CancelAtPeriodEnd
	if transition.Status == models.SubscriptionStatusCanceled || transition.Status == models.SubscriptionStatusExpired {
		now := s.now()
		subscription.CanceledAt = &now
	}
	if transition.ProviderAuthorization != "" {
		if s.secrets == nil {
			return fmt.Errorf("billing secret protector is not configured")
		}
		encrypted, err := s.secrets.Seal(transition.ProviderAuthorization)
		if err != nil {
			return fmt.Errorf("encrypt provider authorization: %w", err)
		}
		subscription.ProviderAuthCode = encrypted
	}
	if err := s.billing.SaveSubscription(ctx, &subscription); err != nil {
		return fmt.Errorf("save billing transition: %w", err)
	}
	return nil
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
