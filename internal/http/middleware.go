package http

import (
	"context"
	"crypto/subtle"
	"fmt"
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

func sessionRequired(auth *services.AuthService, apiKeys *services.APIKeyService, cookieName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		raw := strings.TrimSpace(c.Cookies(cookieName))
		if raw != "" {
			session, err := auth.AuthenticateSession(c.UserContext(), raw)
			if err != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
			}
			c.Locals("session", session)
			return c.Next()
		}
		credential, err := apiKeyFromRequest(c, apiKeys)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		c.Locals("apiKeyUserID", credential.Key.UserID)
		c.Locals("apiKeyCredential", credential)
		return c.Next()
	}
}

func organizationRoleRequired(organizations *services.OrganizationService, required models.MemberRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := authenticatedUserID(c)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "authenticated session is required"})
		}
		if credential, apiKey := c.Locals("apiKeyCredential").(services.APIKeyCredential); apiKey && !apiKeyScopeAllowed(credential, required) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "API key scope is insufficient"})
		}
		if credential, apiKey := c.Locals("apiKeyCredential").(services.APIKeyCredential); apiKey && credential.Key.OrganizationID != nil && *credential.Key.OrganizationID != c.Params("organizationID") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "API key is restricted to another organization"})
		}
		if err := organizations.Authorize(c.UserContext(), c.Params("organizationID"), userID, required); err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Next()
	}
}

func apiKeyFromRequest(c *fiber.Ctx, apiKeys *services.APIKeyService) (services.APIKeyCredential, error) {
	if apiKeys == nil {
		return services.APIKeyCredential{}, fmt.Errorf("authentication is unavailable")
	}
	value := strings.TrimSpace(c.Get("Authorization"))
	parts := strings.SplitN(value, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return services.APIKeyCredential{}, fmt.Errorf("session or bearer API key is required")
	}
	return apiKeys.Authenticate(c.UserContext(), strings.TrimSpace(parts[1]))
}

func authenticatedUserID(c *fiber.Ctx) (string, bool) {
	if session, ok := c.Locals("session").(models.Session); ok && session.UserID != "" {
		return session.UserID, true
	}
	userID, ok := c.Locals("apiKeyUserID").(string)
	return userID, ok && userID != ""
}

func apiKeyScopeAllowed(credential services.APIKeyCredential, required models.MemberRole) bool {
	scope := "organization:read"
	if required == models.MemberRoleMember {
		scope = "organization:write"
	}
	if required == models.MemberRoleAdmin {
		scope = "organization:admin"
	}
	if required == models.MemberRoleOwner {
		scope = "organization:owner"
	}
	for _, granted := range credential.Scopes {
		if granted == "*" || granted == scope {
			return true
		}
	}
	return false
}

func apiKeyScopeRequired(scope string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		credential, ok := c.Locals("apiKeyCredential").(services.APIKeyCredential)
		if !ok {
			return c.Next()
		}
		for _, granted := range credential.Scopes {
			if granted == "*" || granted == scope {
				return c.Next()
			}
		}
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "API key scope is insufficient"})
	}
}

func platformAdminRequired(auth *services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		session, ok := c.Locals("session").(models.Session)
		if !ok || session.UserID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "authenticated session is required"})
		}
		admin, err := auth.IsPlatformAdmin(c.UserContext(), session.UserID)
		if err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "platform admin access is unavailable"})
		}
		if !admin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "platform admin access is required"})
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
