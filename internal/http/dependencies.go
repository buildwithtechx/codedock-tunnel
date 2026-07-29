package http

import (
	"context"
	"fmt"

	"codedock.run/codedock-tunnel/internal/handlers"
	"codedock.run/codedock-tunnel/internal/services"
)

type Dependencies struct {
	Auth          *services.AuthService
	DeviceLogin   *services.DeviceLoginService
	Organizations *services.OrganizationService
	Tunnels       *services.TunnelService
	Agents        *services.AgentService
	Domains       *services.DomainService
	Usage         *services.UsageService
	Billing       *services.BillingService
	OAuth         *services.OAuthService
	Account       *services.AccountService
	Audit         *services.AuditService
	Ready         func(context.Context) error
	PublicAPIURL  string
}

func (d Dependencies) Validate() error {
	if d.Auth == nil || d.DeviceLogin == nil || d.Organizations == nil || d.Tunnels == nil || d.Agents == nil || d.Domains == nil || d.Usage == nil || d.Billing == nil || d.Account == nil || d.Audit == nil {
		return fmt.Errorf("http service dependencies are incomplete")
	}
	return nil
}

type Handlers struct {
	Health              *handlers.HealthHandler
	Auth                *handlers.AuthHandler
	Organizations       *handlers.OrganizationHandler
	Tunnels             *handlers.TunnelHandler
	Agents              *handlers.AgentHandler
	Domains             *handlers.DomainHandler
	Usage               *handlers.UsageHandler
	Billing             *handlers.BillingHandler
	OAuth               *handlers.OAuthHandler
	Account             *handlers.AccountHandler
	auditService        *services.AuditService
	authService         *services.AuthService
	organizationService *services.OrganizationService
}

func buildHandlers(deps Dependencies, cookieName string, cookieSecure bool) (Handlers, error) {
	if err := deps.Validate(); err != nil {
		return Handlers{}, err
	}
	authHandler, err := handlers.NewAuthHandler(deps.Auth, deps.DeviceLogin, cookieName, cookieSecure)
	if err != nil {
		return Handlers{}, err
	}
	organizationHandler, err := handlers.NewOrganizationHandler(deps.Organizations)
	if err != nil {
		return Handlers{}, err
	}
	tunnelHandler, err := handlers.NewTunnelHandler(deps.Tunnels)
	if err != nil {
		return Handlers{}, err
	}
	agentHandler, err := handlers.NewAgentHandler(deps.Agents)
	if err != nil {
		return Handlers{}, err
	}
	domainHandler, err := handlers.NewDomainHandler(deps.Domains)
	if err != nil {
		return Handlers{}, err
	}
	usageHandler, err := handlers.NewUsageHandler(deps.Usage)
	if err != nil {
		return Handlers{}, err
	}
	billingHandler, err := handlers.NewBillingHandler(deps.Billing)
	if err != nil {
		return Handlers{}, err
	}
	var oauthHandler *handlers.OAuthHandler
	if deps.OAuth != nil {
		oauthHandler, err = handlers.NewOAuthHandler(deps.OAuth, deps.PublicAPIURL, cookieName, cookieSecure)
		if err != nil {
			return Handlers{}, err
		}
	}
	accountHandler, err := handlers.NewAccountHandler(deps.Account)
	if err != nil {
		return Handlers{}, err
	}
	return Handlers{Health: handlers.NewHealthHandler(deps.Ready), Auth: authHandler, Organizations: organizationHandler, Tunnels: tunnelHandler, Agents: agentHandler, Domains: domainHandler, Usage: usageHandler, Billing: billingHandler, OAuth: oauthHandler, Account: accountHandler, authService: deps.Auth, organizationService: deps.Organizations, auditService: deps.Audit}, nil
}
