package http

import (
	"context"
	"crypto/subtle"
	"strings"
	"sync"
	"time"

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
		provided := []byte(c.Get("X-Internal-Secret"))
		if expected == "" || subtle.ConstantTimeCompare(provided, []byte(expected)) != 1 {
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

func requestRateLimit(max int, window time.Duration) fiber.Handler {
	type bucket struct {
		started time.Time
		count   int
	}
	var mu sync.Mutex
	buckets := make(map[string]bucket)
	return func(c *fiber.Ctx) error {
		key := c.IP()
		now := time.Now()
		mu.Lock()
		value := buckets[key]
		if value.started.IsZero() || now.Sub(value.started) >= window {
			value = bucket{started: now}
		}
		value.count++
		buckets[key] = value
		allowed := value.count <= max
		mu.Unlock()
		if !allowed {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{"error": "rate limit exceeded"})
		}
		return c.Next()
	}
}
