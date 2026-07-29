package http

import (
	"fmt"
	"strings"

	"codedock.run/codedock-tunnel/internal/models"
	"codedock.run/codedock-tunnel/internal/services"
	"github.com/gofiber/fiber/v2"
)

func sessionRequired(auth *services.AuthService, cookieName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		raw := strings.TrimSpace(c.Cookies(cookieName))
		session, err := auth.AuthenticateSession(c.UserContext(), raw)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		c.Locals("session", session)
		return c.Next()
	}
}

func organizationRoleRequired(organizations *services.OrganizationService, required models.MemberRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		session, ok := c.Locals("session").(models.Session)
		if !ok || session.UserID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "authenticated session is required"})
		}
		if err := organizations.Authorize(c.UserContext(), c.Params("organizationID"), session.UserID, required); err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Next()
	}
}

func internalSecretRequired(expected string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if expected == "" || c.Get("X-Internal-Secret") != expected {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": fmt.Sprintf("invalid internal secret")})
		}
		return c.Next()
	}
}
