package postgres

import (
	"fmt"

	"codedock.run/codedock-tunnel/internal/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("postgres database is required")
	}

	err := db.AutoMigrate(
		&models.User{},
		&models.Session{},
		&models.DeviceLogin{},
		&models.APIKey{},
		&models.Organization{},
		&models.OrganizationMember{},
		&models.Agent{},
		&models.Tunnel{},
		&models.TunnelToken{},
		&models.Domain{},
		&models.Plan{},
		&models.Subscription{},
		&models.BillingEvent{},
		&models.UsageSnapshot{},
		&models.UsageEvent{},
		&models.AuditEvent{},
	)
	if err != nil {
		return fmt.Errorf("migrate postgres models: %w", err)
	}

	return nil
}
