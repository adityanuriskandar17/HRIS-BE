package repository

import (
	"context"

	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	"github.com/google/uuid"
)

// UnitRepository defines the interface for unit data operations
type UnitRepository interface {
	FindAll(ctx context.Context) ([]model.Unit, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Unit, error)
	Create(ctx context.Context, unit model.Unit) (*model.Unit, error)
	Update(ctx context.Context, unit model.Unit) (*model.Unit, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// PositionRepository defines the interface for position data operations
type PositionRepository interface {
	FindAll(ctx context.Context) ([]model.Position, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Position, error)
	Create(ctx context.Context, position model.Position) (*model.Position, error)
	Update(ctx context.Context, position model.Position) (*model.Position, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// EmployeeRepository defines the interface for employee data operations
type EmployeeRepository interface {
	FindAll(ctx context.Context) ([]model.Employee, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Employee, error)
	Create(ctx context.Context, employee model.Employee) (*model.Employee, error)
	Update(ctx context.Context, employee model.Employee) (*model.Employee, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// UserAccountRepository defines the interface for user account data operations
type UserAccountRepository interface {
	FindByEmail(ctx context.Context, email string) (*model.UserAccount, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.UserAccount, error)
	Create(ctx context.Context, user model.UserAccount) (*model.UserAccount, error)
	Update(ctx context.Context, user model.UserAccount) (*model.UserAccount, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// SubscriptionRepository defines the interface for subscription data operations
type SubscriptionRepository interface {
	FindAll(ctx context.Context) ([]model.Subscription, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error)
	FindByTenantID(ctx context.Context, tenantID uuid.UUID) ([]model.Subscription, error)
	Create(ctx context.Context, subscription model.Subscription) (*model.Subscription, error)
	Update(ctx context.Context, subscription model.Subscription) (*model.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// PlanRepository defines the interface for plan data operations
type PlanRepository interface {
	FindAll(ctx context.Context) ([]model.Plan, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Plan, error)
	Create(ctx context.Context, plan model.Plan) (*model.Plan, error)
	Update(ctx context.Context, plan model.Plan) (*model.Plan, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// InvoiceRepository defines the interface for invoice data operations
type InvoiceRepository interface {
	FindAll(ctx context.Context, status string, page, limit int) ([]model.Invoice, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Invoice, error)
	FindBySubscriptionID(ctx context.Context, subscriptionID uuid.UUID) ([]model.Invoice, error)
	Create(ctx context.Context, invoice model.Invoice) (*model.Invoice, error)
	Update(ctx context.Context, invoice model.Invoice) (*model.Invoice, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
