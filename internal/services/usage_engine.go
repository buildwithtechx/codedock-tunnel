package services

import (
	"context"

	"codedock.run/codedock-tunnel/internal/engine"
	"codedock.run/codedock-tunnel/internal/models"
)

type EngineUsageRecorder struct {
	usage *UsageService
}

func NewEngineUsageRecorder(usage *UsageService) (*EngineUsageRecorder, error) {
	if usage == nil {
		return nil, ErrUsageServiceRequired
	}
	return &EngineUsageRecorder{usage: usage}, nil
}

func (r *EngineUsageRecorder) Record(ctx context.Context, measurement engine.UsageMeasurement) error {
	return r.usage.Record(ctx, &models.UsageEvent{OrganizationID: measurement.OrganizationID, TunnelID: stringPointer(measurement.TunnelID), EventType: measurement.EventType, Bytes: measurement.Bytes, Connections: measurement.Connections})
}

func stringPointer(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
