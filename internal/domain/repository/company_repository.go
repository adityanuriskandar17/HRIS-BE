package repository

import (
	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CompanyRepository defines the interface for company data operations
type CompanyRepository interface {
	GetByID(id uuid.UUID) (*model.Company, error)
	Update(id uuid.UUID, company *model.Company) error
}

type CompanyRepositoryImpl struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &CompanyRepositoryImpl{db: db}
}

func (r *CompanyRepositoryImpl) GetByID(id uuid.UUID) (*model.Company, error) {
	var company model.Company
	if err := r.db.Preload("Tenant").First(&company, id).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *CompanyRepositoryImpl) Update(id uuid.UUID, company *model.Company) error {
	return r.db.Model(&model.Company{}).Where("id = ?", id).Updates(company).Error
}