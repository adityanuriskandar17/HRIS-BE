package services

import (
	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/repository"
	"github.com/google/uuid"
)


type tenantService struct {
	tenantRepository repository.TenantRepository
}

func NewTenantService(tenantRepository repository.TenantRepository) *tenantService {
	return &tenantService{tenantRepository: tenantRepository}
}

type TenantService interface {
	GetAll() ([]*model.Tenant, error)
	GetByID(id uuid.UUID) (*model.Tenant, error)
	Create(tenant *model.Tenant) error
	Update(id uuid.UUID, tenant *model.Tenant) error
	Delete(id uuid.UUID) error
}

func (s *tenantService) GetAll() ([]*model.Tenant, error) {
	return s.tenantRepository.GetAll()
}

func (s *tenantService) GetByID(id uuid.UUID) (*model.Tenant, error) {
	return s.tenantRepository.GetByID(id)
}

func (s *tenantService) Create(tenant *model.Tenant) error {
	return s.tenantRepository.Create(tenant)
}

func (s *tenantService) Update(id uuid.UUID, tenant *model.Tenant) error {
	return s.tenantRepository.Update(id, tenant)
}

func (s *tenantService) Delete(id uuid.UUID) error {
	return s.tenantRepository.Delete(id)
}
