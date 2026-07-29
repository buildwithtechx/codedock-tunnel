package http

import (
	"context"
	"strings"

	"codedock.run/codedock-tunnel/internal/models"
	"codedock.run/codedock-tunnel/internal/services"
	"github.com/gofiber/fiber/v2"
)

func auditRequest(audit *services.AuditService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		session, ok := c.Locals("session").(models.Session)
		if audit != nil && ok {
			organizationID := c.Params("organizationID")
			var organization *string
			if organizationID != "" {
				organization = &organizationID
			}
			userID := session.UserID
			_ = audit.Record(context.Background(), &models.AuditEvent{OrganizationID: organization, UserID: &userID, Action: c.Method() + " " + c.Path(), ResourceType: "http", ResourceID: c.Params("tunnelID"), IPAddress: c.IP(), UserAgent: c.Get("User-Agent"), Metadata: `{}`})
		}
		return err
	}
}

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
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid internal secret"})
		}
		return c.Next()
	}
}

func securityHeadersMiddleware(requireTLS bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		if requireTLS || c.Protocol() == "https" || c.Get("X-Forwarded-Proto") == "https" {
			c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		return c.Next()
	}
}
