package http

import (
	"fmt"

	"codedock.run/codedock-tunnel/internal/models"
	"github.com/gofiber/fiber/v2"
)

type RouterOptions struct {
	CookieName        string
	CookieSecure      bool
	InternalAPISecret string
}

func RegisterRoutes(app *fiber.App, handlers Handlers, options RouterOptions) error {
	if app == nil {
		return fmt.Errorf("fiber app is required")
	}
	if handlers.Health == nil || handlers.Auth == nil {
		return fmt.Errorf("health and auth handlers are required")
	}
	app.Get("/healthz", handlers.Health.Liveness)
	app.Get("/readyz", handlers.Health.Readiness)
	app.Post("/api/v1/auth/device/start", handlers.Auth.StartDeviceLogin)
	app.Get("/api/v1/auth/session", handlers.Auth.Session)
	app.Post("/api/v1/auth/logout", handlers.Auth.Logout)

	protected := app.Group("/api/v1", sessionRequired(handlers.authService, options.CookieName))
	protected.Post("/auth/device/complete", handlers.Auth.CompleteDeviceLogin)
	protected.Post("/organizations", handlers.Organizations.Create)
	protected.Post("/organizations/:organizationID/members", organizationRoleRequired(handlers.organizationService, models.MemberRoleAdmin), handlers.Organizations.AddMember)
	protected.Post("/organizations/:organizationID/tunnels", organizationRoleRequired(handlers.organizationService, models.MemberRoleMember), handlers.Tunnels.Create)
	protected.Post("/organizations/:organizationID/agents", organizationRoleRequired(handlers.organizationService, models.MemberRoleAdmin), handlers.Agents.Register)
	protected.Post("/organizations/:organizationID/domains", organizationRoleRequired(handlers.organizationService, models.MemberRoleAdmin), handlers.Domains.Create)
	protected.Get("/organizations/:organizationID/usage/events", organizationRoleRequired(handlers.organizationService, models.MemberRoleViewer), handlers.Usage.Events)
	protected.Get("/organizations/:organizationID/usage/snapshot", organizationRoleRequired(handlers.organizationService, models.MemberRoleViewer), handlers.Usage.Snapshot)
	protected.Get("/organizations/:organizationID/billing", organizationRoleRequired(handlers.organizationService, models.MemberRoleViewer), handlers.Billing.Status)
	protected.Patch("/tunnels/:tunnelID/status", handlers.Tunnels.SetStatus)
	protected.Delete("/tunnels/:tunnelID", handlers.Tunnels.Revoke)
	protected.Post("/domains/:domainID/verify", handlers.Domains.Verify)
	protected.Post("/agents/:agentID/heartbeat", handlers.Agents.Heartbeat)
	protected.Delete("/agents/:agentID", handlers.Agents.Revoke)

	if options.InternalAPISecret != "" {
		app.Get("/internal/health", internalSecretRequired(options.InternalAPISecret), handlers.Health.Readiness)
	}
	return nil
}
