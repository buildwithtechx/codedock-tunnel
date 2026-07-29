package http

import (
	"context"
	"fmt"

	"codedock.run/codedock-tunnel/internal/config"
	"codedock.run/codedock-tunnel/internal/repositories"
	"codedock.run/codedock-tunnel/internal/services"
	"gorm.io/gorm"
)

func NewDatabaseDependencies(db *gorm.DB, cfg config.AuthConfig) (Dependencies, error) {
	if db == nil {
		return Dependencies{}, fmt.Errorf("database is required")
	}
	users, err := repositories.NewUserRepository(db)
	if err != nil {
		return Dependencies{}, err
	}
	identities, err := repositories.NewOAuthIdentityRepository(db)
	if err != nil {
		return Dependencies{}, err
	}
	sessions, err := repositories.NewSessionRepository(db)
	if err != nil {
		return Dependencies{}, err
	}
	deviceLogins, err := repositories.NewDeviceLoginRepository(db)
	if err != nil {
		return Dependencies{}, err
	}
	organizations, err := repositories.NewOrganizationRepository(db)
	if err != nil {
		return Dependencies{}, err
	}
	tunnels, err := repositories.NewTunnelRepository(db)
	if err != nil {
		return Dependencies{}, err
	}
	agents, err := repositories.NewAgentRepository(db)
	if err != nil {
		return Dependencies{}, err
	}
	domains, err := repositories.NewDomainRepository(db)
	if err != nil {
		return Dependencies{}, err
	}
	authService, err := services.NewAuthService(users, identities, sessions, cfg.SessionTTL)
	if err != nil {
		return Dependencies{}, err
	}
	deviceService, err := services.NewDeviceLoginService(deviceLogins, cfg.DeviceLoginTTL)
	if err != nil {
		return Dependencies{}, err
	}
	organizationService, err := services.NewOrganizationService(organizations)
	if err != nil {
		return Dependencies{}, err
	}
	tunnelService, err := services.NewTunnelService(tunnels)
	if err != nil {
		return Dependencies{}, err
	}
	agentService, err := services.NewAgentService(agents)
	if err != nil {
		return Dependencies{}, err
	}
	domainService, err := services.NewDomainService(domains)
	if err != nil {
		return Dependencies{}, err
	}
	return Dependencies{Auth: authService, DeviceLogin: deviceService, Organizations: organizationService, Tunnels: tunnelService, Agents: agentService, Domains: domainService, Ready: func(ctx context.Context) error {
		sqlDB, err := db.DB()
		if err != nil {
			return fmt.Errorf("get database connection: %w", err)
		}
		if err := sqlDB.PingContext(ctx); err != nil {
			return fmt.Errorf("ping database: %w", err)
		}
		return nil
	}}, nil
}
