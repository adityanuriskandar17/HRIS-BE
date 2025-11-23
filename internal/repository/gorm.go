package repository

import (
	"context"

	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UnitRepositoryImpl implements UnitRepository using GORM
type UnitRepositoryImpl struct {
	db *gorm.DB
}

// NewUnitRepository creates a new UnitRepository
func NewUnitRepository(db *gorm.DB) UnitRepository {
	return &UnitRepositoryImpl{db: db}
}

func (r *UnitRepositoryImpl) FindAll(ctx context.Context) ([]model.Unit, error) {
	var units []model.Unit
	if err := r.db.WithContext(ctx).Order("code ASC").Find(&units).Error; err != nil {
		return nil, err
	}
	return units, nil
}

func (r *UnitRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*model.Unit, error) {
	var unit model.Unit
	if err := r.db.WithContext(ctx).First(&unit, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &unit, nil
}

func (r *UnitRepositoryImpl) Create(ctx context.Context, unit model.Unit) (*model.Unit, error) {
	if err := r.db.WithContext(ctx).Create(&unit).Error; err != nil {
		return nil, err
	}
	return &unit, nil
}

func (r *UnitRepositoryImpl) Update(ctx context.Context, unit model.Unit) (*model.Unit, error) {
	if err := r.db.WithContext(ctx).Save(&unit).Error; err != nil {
		return nil, err
	}
	return &unit, nil
}

func (r *UnitRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Unit{}).Error
}

// PositionRepositoryImpl implements PositionRepository using GORM
type PositionRepositoryImpl struct {
	db *gorm.DB
}

// NewPositionRepository creates a new PositionRepository
func NewPositionRepository(db *gorm.DB) PositionRepository {
	return &PositionRepositoryImpl{db: db}
}

func (r *PositionRepositoryImpl) FindAll(ctx context.Context) ([]model.Position, error) {
	var positions []model.Position
	if err := r.db.WithContext(ctx).Order("title ASC").Find(&positions).Error; err != nil {
		return nil, err
	}
	return positions, nil
}

func (r *PositionRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*model.Position, error) {
	var position model.Position
	if err := r.db.WithContext(ctx).First(&position, id).Error; err != nil {
		return nil, err
	}
	return &position, nil
}

func (r *PositionRepositoryImpl) Create(ctx context.Context, position model.Position) (*model.Position, error) {
	if err := r.db.WithContext(ctx).Create(&position).Error; err != nil {
		return nil, err
	}
	return &position, nil
}

func (r *PositionRepositoryImpl) Update(ctx context.Context, position model.Position) (*model.Position, error) {
	if err := r.db.WithContext(ctx).Save(&position).Error; err != nil {
		return nil, err
	}
	return &position, nil
}

func (r *PositionRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Position{}, id).Error
}

// EmployeeRepositoryImpl implements EmployeeRepository using GORM
type EmployeeRepositoryImpl struct {
	db *gorm.DB
}

// NewEmployeeRepository creates a new EmployeeRepository
func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &EmployeeRepositoryImpl{db: db}
}

func (r *EmployeeRepositoryImpl) FindAll(ctx context.Context) ([]model.Employee, error) {
	var employees []model.Employee
	if err := r.db.WithContext(ctx).Preload("Unit").Preload("Position").Order("employee_code ASC").Find(&employees).Error; err != nil {
		return nil, err
	}
	return employees, nil
}

func (r *EmployeeRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*model.Employee, error) {
	var employee model.Employee
	if err := r.db.WithContext(ctx).Preload("Unit").Preload("Position").First(&employee, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *EmployeeRepositoryImpl) Create(ctx context.Context, employee model.Employee) (*model.Employee, error) {
	if err := r.db.WithContext(ctx).Create(&employee).Error; err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *EmployeeRepositoryImpl) Update(ctx context.Context, employee model.Employee) (*model.Employee, error) {
	if err := r.db.WithContext(ctx).Save(&employee).Error; err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *EmployeeRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Employee{}).Error
}

// UserAccountRepositoryImpl implements UserAccountRepository using GORM
type UserAccountRepositoryImpl struct {
	db *gorm.DB
}

// NewUserAccountRepository creates a new UserAccountRepository
func NewUserAccountRepository(db *gorm.DB) UserAccountRepository {
	return &UserAccountRepositoryImpl{db: db}
}

func (r *UserAccountRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.UserAccount, error) {
	var user model.UserAccount
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserAccountRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*model.UserAccount, error) {
	var user model.UserAccount
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserAccountRepositoryImpl) Create(ctx context.Context, user model.UserAccount) (*model.UserAccount, error) {
	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserAccountRepositoryImpl) Update(ctx context.Context, user model.UserAccount) (*model.UserAccount, error) {
	if err := r.db.WithContext(ctx).Save(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserAccountRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.UserAccount{}, id).Error
}

// SubscriptionRepositoryImpl implements SubscriptionRepository using GORM
type SubscriptionRepositoryImpl struct {
	db *gorm.DB
}

// NewSubscriptionRepository creates a new SubscriptionRepository
func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &SubscriptionRepositoryImpl{db: db}
}

func (r *SubscriptionRepositoryImpl) FindAll(ctx context.Context) ([]model.Subscription, error) {
	var subscriptions []model.Subscription
	if err := r.db.WithContext(ctx).Preload("Plan").Preload("Tenant").Order("created_at DESC").Find(&subscriptions).Error; err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func (r *SubscriptionRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	var subscription model.Subscription
	if err := r.db.WithContext(ctx).Preload("Plan").Preload("Tenant").First(&subscription, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (r *SubscriptionRepositoryImpl) FindByTenantID(ctx context.Context, tenantID uuid.UUID) ([]model.Subscription, error) {
	var subscriptions []model.Subscription
	if err := r.db.WithContext(ctx).Preload("Plan").Where("tenant_id = ?", tenantID).Order("created_at DESC").Find(&subscriptions).Error; err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func (r *SubscriptionRepositoryImpl) Create(ctx context.Context, subscription model.Subscription) (*model.Subscription, error) {
	if err := r.db.WithContext(ctx).Create(&subscription).Error; err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (r *SubscriptionRepositoryImpl) Update(ctx context.Context, subscription model.Subscription) (*model.Subscription, error) {
	if err := r.db.WithContext(ctx).Save(&subscription).Error; err != nil {
		return nil, err
	}
	return &subscription, nil
}

func (r *SubscriptionRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Subscription{}, id).Error
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

// InvoiceRepositoryImpl implements InvoiceRepository using GORM
type InvoiceRepositoryImpl struct {
	db *gorm.DB
}

// NewInvoiceRepository creates a new InvoiceRepository
func NewInvoiceRepository(db *gorm.DB) InvoiceRepository {
	return &InvoiceRepositoryImpl{db: db}
}

func (r *InvoiceRepositoryImpl) FindAll(ctx context.Context, status string, page, limit int) ([]model.Invoice, error) {
	var invoices []model.Invoice
	query := r.db.WithContext(ctx).Preload("Subscription.Plan").Preload("Subscription.Tenant")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&invoices).Error; err != nil {
		return nil, err
	}
	return invoices, nil
}

func (r *InvoiceRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*model.Invoice, error) {
	var invoice model.Invoice
	if err := r.db.WithContext(ctx).Preload("Subscription.Plan").Preload("Subscription.Tenant").First(&invoice, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (r *InvoiceRepositoryImpl) FindBySubscriptionID(ctx context.Context, subscriptionID uuid.UUID) ([]model.Invoice, error) {
	var invoices []model.Invoice
	if err := r.db.WithContext(ctx).Preload("Subscription.Plan").Preload("Subscription.Tenant").Where("subscription_id = ?", subscriptionID).Order("created_at DESC").Find(&invoices).Error; err != nil {
		return nil, err
	}
	return invoices, nil
}

func (r *InvoiceRepositoryImpl) Create(ctx context.Context, invoice model.Invoice) (*model.Invoice, error) {
	if err := r.db.WithContext(ctx).Create(&invoice).Error; err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (r *InvoiceRepositoryImpl) Update(ctx context.Context, invoice model.Invoice) (*model.Invoice, error) {
	if err := r.db.WithContext(ctx).Save(&invoice).Error; err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (r *InvoiceRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Invoice{}, id).Error
}
