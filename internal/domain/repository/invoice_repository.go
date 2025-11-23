package repository

import (
	"context"
	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InvoiceRepository defines the interface for invoice data operations
type InvoiceRepository interface {
	FindAll(ctx context.Context, status string, page, limit int) ([]model.Invoice, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Invoice, error)
	FindBySubscriptionID(ctx context.Context, subscriptionID uuid.UUID) ([]model.Invoice, error)
	Create(ctx context.Context, invoice model.Invoice) (*model.Invoice, error)
	Update(ctx context.Context, invoice model.Invoice) (*model.Invoice, error)
	Delete(ctx context.Context, id uuid.UUID) error
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
