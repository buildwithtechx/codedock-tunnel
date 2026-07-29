package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"codedock.run/codedock-tunnel/internal/infra/billing"
	"codedock.run/codedock-tunnel/internal/models"
	"codedock.run/codedock-tunnel/internal/services"
	"github.com/gofiber/fiber/v2"
)

type BillingHandler struct {
	billing       *services.BillingService
	alerts        *services.AlertService
	webhookSecret string
}

type CheckoutRequest struct {
	PlanKey string `json:"planKey" validate:"required"`
}

func NewBillingHandler(billing *services.BillingService) (*BillingHandler, error) {
	if billing == nil {
		return nil, fmt.Errorf("billing service is required")
	}
	return &BillingHandler{billing: billing}, nil
}

func (h *BillingHandler) SetWebhookSecret(secret string)          { h.webhookSecret = secret }
func (h *BillingHandler) SetAlerts(alerts *services.AlertService) { h.alerts = alerts }

func (h *BillingHandler) Status(c *fiber.Ctx) error {
	plan, subscription, err := h.billing.Entitlements(c.UserContext(), c.Params("organizationID"))
	if err != nil {
		return writeError(c, fiber.StatusNotFound, err)
	}
	return c.JSON(fiber.Map{"plan": plan, "subscription": subscription})
}

func (h *BillingHandler) Checkout(c *fiber.Ctx) error {
	var input CheckoutRequest
	if err := c.BodyParser(&input); err != nil {
		return writeError(c, fiber.StatusBadRequest, err)
	}
	url, err := h.billing.Checkout(c.UserContext(), c.Params("organizationID"), input.PlanKey)
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, err)
	}
	return c.JSON(fiber.Map{"url": url})
}

func (h *BillingHandler) Portal(c *fiber.Ctx) error {
	url, err := h.billing.Portal(c.UserContext(), c.Params("organizationID"))
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, err)
	}
	return c.JSON(fiber.Map{"url": url})
}

func (h *BillingHandler) Cancel(c *fiber.Ctx) error {
	if err := h.billing.Cancel(c.UserContext(), c.Params("organizationID")); err != nil {
		return writeError(c, fiber.StatusBadRequest, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *BillingHandler) Resume(c *fiber.Ctx) error {
	if err := h.billing.Resume(c.UserContext(), c.Params("organizationID")); err != nil {
		return writeError(c, fiber.StatusBadRequest, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *BillingHandler) Webhook(c *fiber.Ctx) error {
	payload := c.Body()
	provider := c.Params("provider")
	if h.webhookSecret == "" || !billing.VerifyHMACSHA256(payload, c.Get("X-Signature"), h.webhookSecret) {
		if h.alerts != nil {
			_ = h.alerts.AlertFailedWebhook(c.UserContext(), provider, c.Get("X-Event-ID"), "invalid billing webhook signature")
		}
		return writeError(c, fiber.StatusUnauthorized, fmt.Errorf("invalid billing webhook signature"))
	}
	eventID := c.Get("X-Event-ID")
	if eventID == "" {
		if h.alerts != nil {
			_ = h.alerts.AlertFailedWebhook(c.UserContext(), provider, "", "missing X-Event-ID header")
		}
		return writeError(c, fiber.StatusBadRequest, fmt.Errorf("X-Event-ID is required"))
	}
	digest := sha256.Sum256(payload)
	event := &models.BillingEvent{Provider: models.BillingProvider(provider), ProviderEventID: eventID, EventType: c.Get("X-Event-Type"), PayloadHash: hex.EncodeToString(digest[:])}
	created, err := h.billing.RecordEvent(c.UserContext(), event)
	if err != nil {
		if h.alerts != nil {
			_ = h.alerts.AlertFailedWebhook(c.UserContext(), provider, eventID, err.Error())
		}
		return writeError(c, fiber.StatusBadRequest, err)
	}
	if created {
		if err := h.billing.MarkProcessed(c.UserContext(), event.ID); err != nil {
			if h.alerts != nil {
				_ = h.alerts.AlertFailedWebhook(c.UserContext(), provider, eventID, err.Error())
			}
			return writeError(c, fiber.StatusInternalServerError, err)
		}
	}
	return c.SendStatus(fiber.StatusAccepted)
}
