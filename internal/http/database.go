package http

import (
	"context"
	"fmt"

	"codedock.run/codedock-tunnel/internal/config"
	"codedock.run/codedock-tunnel/internal/repositories"
	"codedock.run/codedock-tunnel/internal/services"
	"gorm.io/gorm"
)

func NewDatabaseDependencies(db *gorm.DB, cfg config.APIConfig) (Dependencies, error) {
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
	billingRepository, err := repositories.NewBillingRepository(db)
	if err != nil {
		return Dependencies{}, err
	}
	authService, err := services.NewAuthService(users, identities, sessions, cfg.Auth.SessionTTL)
	if err != nil {
		return Dependencies{}, err
	}
	deviceService, err := services.NewDeviceLoginService(deviceLogins, cfg.Auth.DeviceLoginTTL)
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
	billingService, err := services.NewBillingService(billingRepository)
	if err != nil {
		return Dependencies{}, err
	}
	domainService, err := services.NewDomainService(domains)
	if err != nil {
		return Dependencies{}, err
	}
	tunnelService.SetBilling(billingService)
	organizationService.SetBilling(billingService)
	domainService.SetBilling(billingService)
	hostnameAllocator, err := services.NewHostnameAllocator(tunnels, cfg.Tunnel.Domain)
	if err != nil {
		return Dependencies{}, err
	}
	tunnelService.SetHostnameAllocator(hostnameAllocator)
	agentService, err := services.NewAgentService(agents)
	if err != nil {
		return Dependencies{}, err
	}
	usageRepository, err := repositories.NewUsageRepository(db)
	if err != nil {
		return Dependencies{}, err
	}
	usageService, err := services.NewUsageService(usageRepository)
	if err != nil {
		return Dependencies{}, err
	}
	auditRepository, err := repositories.NewAuditRepository(db)
	if err != nil {
		return Dependencies{}, err
	}
	auditService, err := services.NewAuditService(auditRepository)
	if err != nil {
		return Dependencies{}, err
	}
	accountService, err := services.NewAccountService(users, organizations)
	if err != nil {
		return Dependencies{}, err
	}
	return Dependencies{Auth: authService, DeviceLogin: deviceService, Organizations: organizationService, Tunnels: tunnelService, Agents: agentService, Domains: domainService, Usage: usageService, Billing: billingService, Account: accountService, Audit: auditService, Ready: func(ctx context.Context) error {
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
