package repository

import (
	"context"
	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PlanRepository defines the interface for plan data operations
type PlanRepository interface {
	FindAll(ctx context.Context) ([]model.Plan, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Plan, error)
	Create(ctx context.Context, plan model.Plan) (*model.Plan, error)
	Update(ctx context.Context, plan model.Plan) (*model.Plan, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// PlanRepositoryImpl implements PlanRepository using GORM
type PlanRepositoryImpl struct {
	db *gorm.DB
}

// NewPlanRepository creates a new PlanRepository
func NewPlanRepository(db *gorm.DB) PlanRepository {
	return &PlanRepositoryImpl{db: db}
}

func (r *PlanRepositoryImpl) FindAll(ctx context.Context) ([]model.Plan, error) {
	var plans []model.Plan
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&plans).Error; err != nil {
		return nil, err
	}
	return plans, nil
}

func (r *PlanRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*model.Plan, error) {
	var plan model.Plan
	if err := r.db.WithContext(ctx).First(&plan, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *PlanRepositoryImpl) Create(ctx context.Context, plan model.Plan) (*model.Plan, error) {
	if err := r.db.WithContext(ctx).Create(&plan).Error; err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *PlanRepositoryImpl) Update(ctx context.Context, plan model.Plan) (*model.Plan, error) {
	if err := r.db.WithContext(ctx).Save(&plan).Error; err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *PlanRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Plan{}, id).Error
}