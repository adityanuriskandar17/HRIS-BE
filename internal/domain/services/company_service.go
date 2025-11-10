package services

import (
	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/repository"
	"github.com/google/uuid"
)

type companyService struct {
	companyRepository repository.CompanyRepository
}

func NewCompanyService(companyRepository repository.CompanyRepository) CompanyService {
	return &companyService{companyRepository: companyRepository}
}

type CompanyService interface {
	GetByID(id uuid.UUID) (*model.Company, error)
	Update(id uuid.UUID, company *model.Company) error
	GetLimits(id uuid.UUID) (*model.Company, error)
}

func (s *companyService) GetByID(id uuid.UUID) (*model.Company, error) {
	return s.companyRepository.GetByID(id)
}

func (s *companyService) Update(id uuid.UUID, company *model.Company) error {
	return s.companyRepository.Update(id, company)
}

func (s *companyService) GetLimits(id uuid.UUID) (*model.Company, error) {
	return s.companyRepository.GetByID(id)
}