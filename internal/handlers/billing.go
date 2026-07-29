package handlers

import (
	"fmt"

	"codedock.run/codedock-tunnel/internal/services"
	"github.com/gofiber/fiber/v2"
)

type BillingHandler struct{ billing *services.BillingService }

func NewBillingHandler(billing *services.BillingService) (*BillingHandler, error) {
	if billing == nil {
		return nil, fmt.Errorf("billing service is required")
	}
	return &BillingHandler{billing: billing}, nil
}

func (h *BillingHandler) Status(c *fiber.Ctx) error {
	plan, subscription, err := h.billing.Entitlements(c.UserContext(), c.Params("organizationID"))
	if err != nil {
		return writeError(c, fiber.StatusNotFound, err)
	}
	return c.JSON(fiber.Map{"plan": plan, "subscription": subscription})
}
