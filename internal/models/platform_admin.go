package models

type PlatformAdminRole string

const (
	PlatformAdminOwner PlatformAdminRole = "owner"
	PlatformAdminAdmin PlatformAdminRole = "admin"
)

type PlatformAdmin struct {
	Base
	UserID string            `json:"userId" gorm:"type:uuid;uniqueIndex;not null"`
	Name   string            `json:"name" gorm:"type:varchar(120);not null;default:''"`
	Role   PlatformAdminRole `json:"role" gorm:"type:varchar(20);not null;default:'admin';index"`
	Active bool              `json:"active" gorm:"not null;default:true;index"`
}
