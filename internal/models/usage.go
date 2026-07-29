package models

import "time"

type UsageSnapshot struct {
	Base
	OrganizationID    string    `json:"organizationId" gorm:"type:uuid;not null;uniqueIndex:organization_period"`
	PeriodStart       time.Time `json:"periodStart" gorm:"not null;uniqueIndex:organization_period"`
	PeriodEnd         time.Time `json:"periodEnd" gorm:"not null"`
	TunnelCount       int       `json:"tunnelCount" gorm:"not null;default:0"`
	ActiveConnections int       `json:"activeConnections" gorm:"not null;default:0"`
	BandwidthBytes    int64     `json:"bandwidthBytes" gorm:"not null;default:0"`
	RequestCount      int64     `json:"requestCount" gorm:"not null;default:0"`
	ErrorCount        int64     `json:"errorCount" gorm:"not null;default:0"`
}

type UsageEvent struct {
	Base
	OrganizationID string    `json:"organizationId" gorm:"type:uuid;not null;index"`
	TunnelID       *string   `json:"tunnelId,omitempty" gorm:"type:uuid;index"`
	EventType      string    `json:"eventType" gorm:"not null;index"`
	Bytes          int64     `json:"bytes" gorm:"not null;default:0"`
	Connections    int       `json:"connections" gorm:"not null;default:0"`
	OccurredAt     time.Time `json:"occurredAt" gorm:"not null;index"`
}
