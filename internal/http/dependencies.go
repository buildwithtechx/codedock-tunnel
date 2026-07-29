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
	Ready         func(context.Context) error
}

func (d Dependencies) Validate() error {
	if d.Auth == nil || d.DeviceLogin == nil || d.Organizations == nil || d.Tunnels == nil || d.Agents == nil || d.Domains == nil {
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
	return Handlers{Health: handlers.NewHealthHandler(deps.Ready), Auth: authHandler, Organizations: organizationHandler, Tunnels: tunnelHandler, Agents: agentHandler, Domains: domainHandler, authService: deps.Auth, organizationService: deps.Organizations}, nil
}
