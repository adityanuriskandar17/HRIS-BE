package repository

import (
	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TenantRepository defines the interface for tenant data operations
type TenantRepository interface {
	GetAll() ([]*model.Tenant, error)
	GetByID(id uuid.UUID) (*model.Tenant, error)
	Create(tenant *model.Tenant) error
	Update(id uuid.UUID, tenant *model.Tenant) error
	Delete(id uuid.UUID) error
}

type TenantRepositoryImpl struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) TenantRepository {
	return &TenantRepositoryImpl{db: db}
}

func (r *TenantRepositoryImpl) GetAll() ([]*model.Tenant, error) {
	var tenants []*model.Tenant
	if err := r.db.Find(&tenants).Error; err != nil {
		return nil, err
	}
	return tenants, nil
}

func (r *TenantRepositoryImpl) GetByID(id uuid.UUID) (*model.Tenant, error) {
	var tenant model.Tenant
	if err := r.db.First(&tenant, id).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *TenantRepositoryImpl) Create(tenant *model.Tenant) error {
	return r.db.Create(tenant).Error
}

func (r *TenantRepositoryImpl) Update(id uuid.UUID, tenant *model.Tenant) error {
	return r.db.Model(&model.Tenant{}).Where("id = ?", id).Updates(tenant).Error
}

func (r *TenantRepositoryImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Tenant{}, id).Error
}
