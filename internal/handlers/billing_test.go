package handlers

import (
	"testing"

	"codedock.run/codedock-tunnel/internal/models"
)

func TestParseWebhookCapturesPaymentFailureData(t *testing.T) {
	transition, eventID, eventType, err := parseWebhook(string(models.BillingProviderPaystack), []byte(`{"id":"evt_1","event":"charge.failed","data":{"subscription_code":"sub_1","amount":2500,"currency":"NGN","attempts_remaining":2}}`))
	if err != nil {
		t.Fatalf("parse webhook: %v", err)
	}
	if eventID != "evt_1" || eventType != "charge.failed" {
		t.Fatalf("unexpected event metadata: %s %s", eventID, eventType)
	}
	if transition.Status != models.SubscriptionStatusPastDue || transition.AmountMinor != 2500 || transition.AttemptsRemaining != 2 {
		t.Fatalf("unexpected payment transition: %+v", transition)
	}
}

func TestParseWebhookCapturesSubscriptionResetData(t *testing.T) {
	transition, _, eventType, err := parseWebhook(string(models.BillingProviderPolar), []byte(`{"id":"evt_2","type":"subscription.downgraded","data":{"subscription_id":"sub_2","previous_plan":"Beam"}}`))
	if err != nil {
		t.Fatalf("parse webhook: %v", err)
	}
	if eventType != "subscription.downgraded" || transition.PreviousPlan != "Beam" {
		t.Fatalf("unexpected reset transition: %+v", transition)
	}
}
