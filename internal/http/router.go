package http

import (
	"fmt"
	"time"

	"codedock.run/codedock-tunnel/internal/models"
	"github.com/gofiber/fiber/v2"
)

type RouterOptions struct {
	CookieName           string
	CookieSecure         bool
	InternalAPISecret    string
	BillingWebhookSecret string
}

func RegisterRoutes(app *fiber.App, handlers Handlers, options RouterOptions) error {
	if app == nil {
		return fmt.Errorf("fiber app is required")
	}
	if handlers.Health == nil || handlers.Auth == nil {
		return fmt.Errorf("health and auth handlers are required")
	}
	if handlers.Billing != nil {
		handlers.Billing.SetWebhookSecret(options.BillingWebhookSecret)
	}
	app.Use(securityHeadersMiddleware(options.CookieSecure))
	app.Get("/healthz", handlers.Health.Liveness)
	app.Get("/readyz", handlers.Health.Readiness)
	app.Get("/metrics", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/plain; version=0.0.4")
		return c.SendString("# HELP codedock_status Status metric\n# TYPE codedock_status gauge\ncodedock_status 1\n")
	})
	authLimiter := requestRateLimit(10, time.Minute)
	app.Post("/api/v1/auth/device/start", authLimiter, handlers.Auth.StartDeviceLogin)
	app.Get("/api/v1/auth/device/poll", authLimiter, handlers.Auth.PollDeviceLogin)
	if handlers.Billing != nil {
		app.Post("/api/v1/billing/webhooks/:provider", handlers.Billing.Webhook)
	}
	if handlers.OAuth != nil {
		app.Get("/api/v1/auth/oauth/:provider", handlers.OAuth.Start)
		app.Get("/api/v1/auth/oauth/:provider/callback", handlers.OAuth.Callback)
	}
	app.Get("/api/v1/auth/session", handlers.Auth.Session)
	app.Post("/api/v1/auth/logout", handlers.Auth.Logout)

	protected := app.Group("/api/v1", sessionRequired(handlers.authService, options.CookieName), auditRequest(handlers.auditService))
	protected.Post("/auth/device/complete", handlers.Auth.CompleteDeviceLogin)
	protected.Post("/organizations", handlers.Organizations.Create)
	protected.Post("/organizations/:organizationID/members", organizationRoleRequired(handlers.organizationService, models.MemberRoleAdmin), handlers.Organizations.AddMember)
	protected.Post("/organizations/:organizationID/invitations", organizationRoleRequired(handlers.organizationService, models.MemberRoleAdmin), handlers.Invitations.Create)
	protected.Post("/invitations/accept", handlers.Invitations.Accept)
	protected.Delete("/account", handlers.Account.Delete)
	protected.Post("/organizations/:organizationID/transfer", organizationRoleRequired(handlers.organizationService, models.MemberRoleOwner), handlers.Account.TransferOwnership)
	protected.Post("/organizations/:organizationID/tunnels", organizationRoleRequired(handlers.organizationService, models.MemberRoleMember), handlers.Tunnels.Create)
	protected.Get("/organizations/:organizationID/tunnels", organizationRoleRequired(handlers.organizationService, models.MemberRoleViewer), handlers.Tunnels.List)
	protected.Post("/organizations/:organizationID/agents", organizationRoleRequired(handlers.organizationService, models.MemberRoleAdmin), handlers.Agents.Register)
	protected.Post("/organizations/:organizationID/domains", organizationRoleRequired(handlers.organizationService, models.MemberRoleAdmin), handlers.Domains.Create)
	protected.Get("/organizations/:organizationID/usage/events", organizationRoleRequired(handlers.organizationService, models.MemberRoleViewer), handlers.Usage.Events)
	protected.Get("/organizations/:organizationID/usage/snapshot", organizationRoleRequired(handlers.organizationService, models.MemberRoleViewer), handlers.Usage.Snapshot)
	protected.Get("/organizations/:organizationID/billing", organizationRoleRequired(handlers.organizationService, models.MemberRoleViewer), handlers.Billing.Status)
	protected.Post("/organizations/:organizationID/billing/checkout", organizationRoleRequired(handlers.organizationService, models.MemberRoleOwner), handlers.Billing.Checkout)
	protected.Get("/organizations/:organizationID/billing/portal", organizationRoleRequired(handlers.organizationService, models.MemberRoleOwner), handlers.Billing.Portal)
	protected.Post("/organizations/:organizationID/billing/cancel", organizationRoleRequired(handlers.organizationService, models.MemberRoleOwner), handlers.Billing.Cancel)
	protected.Post("/organizations/:organizationID/billing/resume", organizationRoleRequired(handlers.organizationService, models.MemberRoleOwner), handlers.Billing.Resume)
	protected.Patch("/tunnels/:tunnelID/status", handlers.Tunnels.SetStatus)
	protected.Get("/tunnels/:tunnelID", handlers.Tunnels.Inspect)
	protected.Delete("/tunnels/:tunnelID", handlers.Tunnels.Revoke)
	protected.Post("/domains/:domainID/verify", handlers.Domains.Verify)
	protected.Post("/agents/:agentID/heartbeat", handlers.Agents.Heartbeat)
	protected.Delete("/agents/:agentID", handlers.Agents.Revoke)

	admin := protected.Group("/admin", platformAdminRequired(handlers.authService))
	admin.Get("/overview", handlers.Admin.Overview)
	admin.Get("/users", handlers.Admin.Users)
	admin.Patch("/users/:userID/status", handlers.Admin.SetUserStatus)
	admin.Get("/organizations", handlers.Admin.Organizations)
	admin.Get("/tunnels", handlers.Admin.Tunnels)
	admin.Get("/subscriptions", handlers.Admin.Subscriptions)
	admin.Get("/usage", handlers.Admin.Usage)
	admin.Get("/audit-logs", handlers.Admin.AuditLogs)
	admin.Post("/tunnels/:tunnelID/revoke", handlers.Tunnels.Revoke)

	if options.InternalAPISecret != "" {
		app.Get("/internal/health", internalSecretRequired(options.InternalAPISecret), handlers.Health.Readiness)
		app.Get("/internal/agents/authenticate", internalSecretRequired(options.InternalAPISecret), handlers.Agents.Authenticate)
		app.Post("/internal/usage", internalSecretRequired(options.InternalAPISecret), handlers.Usage.Ingest)
	}
	return nil
}
