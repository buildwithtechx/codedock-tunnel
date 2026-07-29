package repositories

import (
	"context"
	"fmt"
	"time"

	"codedock.run/codedock-tunnel/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UsageRepository interface {
	CreateEvent(context.Context, *models.UsageEvent) error
	UpsertSnapshot(context.Context, *models.UsageSnapshot) error
	FindSnapshot(context.Context, string, time.Time) (models.UsageSnapshot, error)
	ListEvents(context.Context, string, time.Time, time.Time) ([]models.UsageEvent, error)
}

type GormUsageRepository struct{ db *gorm.DB }

func NewUsageRepository(db *gorm.DB) (*GormUsageRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("database is required")
	}
	return &GormUsageRepository{db: db}, nil
}

func (r *GormUsageRepository) CreateEvent(ctx context.Context, event *models.UsageEvent) error {
	if event == nil {
		return fmt.Errorf("usage event is required")
	}
	return wrap(r.db.WithContext(ctx).Create(event).Error, "create usage event")
}

func (r *GormUsageRepository) UpsertSnapshot(ctx context.Context, snapshot *models.UsageSnapshot) error {
	if snapshot == nil {
		return fmt.Errorf("usage snapshot is required")
	}
	return wrap(r.db.WithContext(ctx).Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "organization_id"}, {Name: "period_start"}}, DoUpdates: clause.AssignmentColumns([]string{"period_end", "tunnel_count", "active_connections", "bandwidth_bytes", "request_count", "error_count", "updated_at"})}).Create(snapshot).Error, "upsert usage snapshot")
}

func (r *GormUsageRepository) FindSnapshot(ctx context.Context, organizationID string, periodStart time.Time) (models.UsageSnapshot, error) {
	var snapshot models.UsageSnapshot
	if err := r.db.WithContext(ctx).Where("organization_id = ? AND period_start = ?", organizationID, periodStart).First(&snapshot).Error; err != nil {
		return models.UsageSnapshot{}, mapError(err)
	}
	return snapshot, nil
}

func (r *GormUsageRepository) ListEvents(ctx context.Context, organizationID string, from, to time.Time) ([]models.UsageEvent, error) {
	var events []models.UsageEvent
	err := r.db.WithContext(ctx).Where("organization_id = ? AND occurred_at >= ? AND occurred_at < ?", organizationID, from, to).Order("occurred_at ASC").Find(&events).Error
	if err != nil {
		return nil, fmt.Errorf("list usage events: %w", err)
	}
	return events, nil
}
