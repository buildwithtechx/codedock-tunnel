package repositories

import (
	"context"
	"fmt"
	"time"

	"codedock.run/codedock-tunnel/internal/models"
	"gorm.io/gorm"
)

type TimeSeriesUsageRepository interface {
	RecordTimeSeriesEvent(ctx context.Context, event models.UsageEvent) error
	AggregateTimeSeries(ctx context.Context, organizationID string, start, end time.Time) (models.UsageSnapshot, error)
}

type GormTimeSeriesUsageRepository struct {
	db *gorm.DB
}

func NewGormTimeSeriesUsageRepository(db *gorm.DB) (*GormTimeSeriesUsageRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is required")
	}
	return &GormTimeSeriesUsageRepository{db: db}, nil
}

func (r *GormTimeSeriesUsageRepository) RecordTimeSeriesEvent(ctx context.Context, event models.UsageEvent) error {
	if event.OrganizationID == "" || event.EventType == "" {
		return fmt.Errorf("organization and event type are required for time series tracking")
	}

	sanitized := models.UsageEvent{
		Base:           event.Base,
		OrganizationID: event.OrganizationID,
		TunnelID:       event.TunnelID,
		EventType:      event.EventType,
		Bytes:          event.Bytes,
		Connections:    event.Connections,
		OccurredAt:     event.OccurredAt,
	}
	if err := r.db.WithContext(ctx).Table("usage_analytics_events").Create(&sanitized).Error; err != nil {
		return fmt.Errorf("record time series usage event: %w", err)
	}
	return nil
}

func (r *GormTimeSeriesUsageRepository) AggregateTimeSeries(ctx context.Context, organizationID string, start, end time.Time) (models.UsageSnapshot, error) {
	var result struct {
		TotalBytes       int64 `gorm:"column:total_bytes"`
		TotalConnections int   `gorm:"column:total_connections"`
		TotalRequests    int64 `gorm:"column:total_requests"`
	}
	query := r.db.WithContext(ctx).Table("usage_analytics_events").
		Select("COALESCE(SUM(bytes), 0) as total_bytes, COALESCE(SUM(connections), 0) as total_connections, COUNT(*) as total_requests").
		Where("organization_id = ? AND occurred_at >= ? AND occurred_at <= ?", organizationID, start, end).
		Scan(&result)
	if query.Error != nil {
		return models.UsageSnapshot{}, fmt.Errorf("aggregate time series usage: %w", query.Error)
	}
	return models.UsageSnapshot{
		OrganizationID:    organizationID,
		PeriodStart:       start,
		PeriodEnd:         end,
		ActiveConnections: result.TotalConnections,
		BandwidthBytes:    result.TotalBytes,
		RequestCount:      result.TotalRequests,
	}, nil
}
