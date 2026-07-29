package models

type Organization struct {
	Base
	Name     string `json:"name" gorm:"not null"`
	Slug     string `json:"slug" gorm:"uniqueIndex;not null"`
	OwnerID  string `json:"ownerId" gorm:"type:uuid;not null;index"`
	Settings string `json:"settings" gorm:"type:jsonb;not null;default:'{}'"`
}

type OrganizationMember struct {
	Base
	OrganizationID string     `json:"organizationId" gorm:"type:uuid;not null;uniqueIndex:organization_user"`
	UserID         string     `json:"userId" gorm:"type:uuid;not null;uniqueIndex:organization_user"`
	Role           MemberRole `json:"role" gorm:"type:varchar(20);not null;default:'member'"`
}
