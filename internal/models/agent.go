package models

import "time"

type AgentStatus string

const (
	AgentStatusPending AgentStatus = "pending"
	AgentStatusOnline  AgentStatus = "online"
	AgentStatusOffline AgentStatus = "offline"
	AgentStatusRevoked AgentStatus = "revoked"
)

type Agent struct {
	Base
	OrganizationID string      `json:"organizationId" gorm:"type:uuid;not null;index"`
	Name           string      `json:"name" gorm:"not null"`
	TokenHash      string      `json:"-" gorm:"uniqueIndex;not null"`
	Status         AgentStatus `json:"status" gorm:"type:varchar(20);not null;default:'pending';index"`
	Version        string      `json:"version,omitempty"`
	Hostname       string      `json:"hostname,omitempty"`
	Platform       string      `json:"platform,omitempty"`
	LastSeenAt     *time.Time  `json:"lastSeenAt,omitempty"`
	ConnectedAt    *time.Time  `json:"connectedAt,omitempty"`
	RevokedAt      *time.Time  `json:"revokedAt,omitempty"`
	Metadata       string      `json:"metadata" gorm:"type:jsonb;not null;default:'{}'"`
}
