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
	GetAll() ([]*model.Company, error)
	GetByID(id uuid.UUID) (*model.Company, error)
	Create(company *model.Company) error
	Update(id uuid.UUID, company *model.Company) error
	GetLimits(id uuid.UUID) (*model.Company, error)
}

func (s *companyService) GetAll() ([]*model.Company, error) {
	return s.companyRepository.GetAll()
}

func (s *companyService) GetByID(id uuid.UUID) (*model.Company, error) {
	return s.companyRepository.GetByID(id)
}

func (s *companyService) Create(company *model.Company) error {
	return s.companyRepository.Create(company)
}

func (s *companyService) Update(id uuid.UUID, company *model.Company) error {
	return s.companyRepository.Update(id, company)
}

func (s *companyService) GetLimits(id uuid.UUID) (*model.Company, error) {
	return s.companyRepository.GetByID(id)
}
