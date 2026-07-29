package repositories

import (
	"context"
	"fmt"

	"codedock.run/codedock-tunnel/internal/models"
	"gorm.io/gorm"
)

type OrganizationRepository interface {
	Create(context.Context, *models.Organization) error
	FindByID(context.Context, string) (models.Organization, error)
	FindBySlug(context.Context, string) (models.Organization, error)
	Update(context.Context, *models.Organization) error
	AddMember(context.Context, *models.OrganizationMember) error
	FindMember(context.Context, string, string) (models.OrganizationMember, error)
	ListMembers(context.Context, string) ([]models.OrganizationMember, error)
	RemoveMember(context.Context, string, string) error
}

type GormOrganizationRepository struct{ db *gorm.DB }

func NewOrganizationRepository(db *gorm.DB) (*GormOrganizationRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("database is required")
	}
	return &GormOrganizationRepository{db: db}, nil
}

func (r *GormOrganizationRepository) Create(ctx context.Context, organization *models.Organization) error {
	if organization == nil {
		return fmt.Errorf("organization is required")
	}
	return wrap(r.db.WithContext(ctx).Create(organization).Error, "create organization")
}

func (r *GormOrganizationRepository) FindByID(ctx context.Context, id string) (models.Organization, error) {
	var organization models.Organization
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&organization).Error
	if err != nil {
		return models.Organization{}, mapError(err)
	}
	return organization, nil
}

func (r *GormOrganizationRepository) FindBySlug(ctx context.Context, slug string) (models.Organization, error) {
	var organization models.Organization
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&organization).Error
	if err != nil {
		return models.Organization{}, mapError(err)
	}
	return organization, nil
}

func (r *GormOrganizationRepository) Update(ctx context.Context, organization *models.Organization) error {
	if organization == nil {
		return fmt.Errorf("organization is required")
	}
	return wrap(r.db.WithContext(ctx).Save(organization).Error, "update organization")
}

func (r *GormOrganizationRepository) AddMember(ctx context.Context, member *models.OrganizationMember) error {
	if member == nil {
		return fmt.Errorf("organization member is required")
	}
	return wrap(r.db.WithContext(ctx).Create(member).Error, "add organization member")
}

func (r *GormOrganizationRepository) FindMember(ctx context.Context, organizationID, userID string) (models.OrganizationMember, error) {
	var member models.OrganizationMember
	err := r.db.WithContext(ctx).Where("organization_id = ? AND user_id = ?", organizationID, userID).First(&member).Error
	if err != nil {
		return models.OrganizationMember{}, mapError(err)
	}
	return member, nil
}

func (r *GormOrganizationRepository) ListMembers(ctx context.Context, organizationID string) ([]models.OrganizationMember, error) {
	var members []models.OrganizationMember
	err := r.db.WithContext(ctx).Where("organization_id = ?", organizationID).Order("created_at ASC").Find(&members).Error
	if err != nil {
		return nil, fmt.Errorf("list organization members: %w", err)
	}
	return members, nil
}

func (r *GormOrganizationRepository) RemoveMember(ctx context.Context, organizationID, userID string) error {
	result := r.db.WithContext(ctx).Where("organization_id = ? AND user_id = ?", organizationID, userID).Delete(&models.OrganizationMember{})
	if result.Error != nil {
		return fmt.Errorf("remove organization member: %w", result.Error)
	}
	if result.RowsAffected != 1 {
		return ErrNotFound
	}
	return nil
}
