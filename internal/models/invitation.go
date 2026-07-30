package models

import "time"

type OrganizationInvitation struct {
	Base
	OrganizationID string     `json:"organizationId" gorm:"type:uuid;not null;index"`
	InviterID      string     `json:"inviterId" gorm:"type:uuid;not null;index"`
	Email          string     `json:"email" gorm:"not null;index"`
	Role           MemberRole `json:"role" gorm:"type:varchar(20);not null"`
	TokenHash      string     `json:"-" gorm:"uniqueIndex;not null"`
	ExpiresAt      time.Time  `json:"expiresAt" gorm:"not null;index"`
	AcceptedAt     *time.Time `json:"acceptedAt,omitempty"`
	RevokedAt      *time.Time `json:"-"`
}
