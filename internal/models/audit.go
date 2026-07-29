package models

import "time"

type AuditEvent struct {
	Base
	OrganizationID *string   `json:"organizationId,omitempty" gorm:"type:uuid;index"`
	UserID         *string   `json:"userId,omitempty" gorm:"type:uuid;index"`
	Action         string    `json:"action" gorm:"not null;index"`
	ResourceType   string    `json:"resourceType" gorm:"not null"`
	ResourceID     string    `json:"resourceId,omitempty" gorm:"index"`
	IPAddress      string    `json:"ipAddress,omitempty"`
	UserAgent      string    `json:"userAgent,omitempty"`
	Metadata       string    `json:"metadata" gorm:"type:jsonb;not null;default:'{}'"`
	OccurredAt     time.Time `json:"occurredAt" gorm:"not null;index"`
}
