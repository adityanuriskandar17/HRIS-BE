package services

import (
	"context"
	"time"

	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/repository"
	"github.com/google/uuid"
)

type invoiceService struct {
	invoiceRepository     repository.InvoiceRepository
	subscriptionRepository repository.SubscriptionRepository
}

func NewInvoiceService(invoiceRepository repository.InvoiceRepository, subscriptionRepository repository.SubscriptionRepository) InvoiceService {
	return &invoiceService{
		invoiceRepository:     invoiceRepository,
		subscriptionRepository: subscriptionRepository,
	}
}

type InvoiceService interface {
	GetAll(ctx context.Context, status string, page, limit int) ([]*model.Invoice, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Invoice, error)
	GetBySubscriptionID(ctx context.Context, subscriptionID uuid.UUID) ([]*model.Invoice, error)
	Create(ctx context.Context, invoice *model.Invoice) error
	Update(ctx context.Context, id uuid.UUID, invoice *model.Invoice) error
	Send(ctx context.Context, id uuid.UUID) error
	Pay(ctx context.Context, id uuid.UUID) error
}

func (s *invoiceService) GetAll(ctx context.Context, status string, page, limit int) ([]*model.Invoice, error) {
	invoices, err := s.invoiceRepository.FindAll(ctx, status, page, limit)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Invoice, len(invoices))
	for i, inv := range invoices {
		result[i] = &inv
	}

	return result, nil
}

func (s *invoiceService) GetByID(ctx context.Context, id uuid.UUID) (*model.Invoice, error) {
	return s.invoiceRepository.FindByID(ctx, id)
}

func (s *invoiceService) GetBySubscriptionID(ctx context.Context, subscriptionID uuid.UUID) ([]*model.Invoice, error) {
	invoices, err := s.invoiceRepository.FindBySubscriptionID(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Invoice, len(invoices))
	for i, inv := range invoices {
		result[i] = &inv
	}

	return result, nil
}

func (s *invoiceService) Create(ctx context.Context, invoice *model.Invoice) error {
	// Set default values
	invoice.ID = uuid.New()
	invoice.Status = model.InvoiceStatusDraft
	invoice.CreatedAt = time.Now()
	invoice.UpdatedAt = time.Now()

	// Verify subscription exists
	_, err := s.subscriptionRepository.FindByID(ctx, invoice.SubscriptionID)
	if err != nil {
		return err
	}

	_, err = s.invoiceRepository.Create(ctx, *invoice)
	return err
}

func (s *invoiceService) Update(ctx context.Context, id uuid.UUID, invoice *model.Invoice) error {
	// Get existing invoice
	existing, err := s.invoiceRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Update fields
	existing.SubscriptionID = invoice.SubscriptionID
	existing.Amount = invoice.Amount
	existing.DueDate = invoice.DueDate
	existing.Status = invoice.Status
	existing.UpdatedAt = time.Now()

	_, err = s.invoiceRepository.Update(ctx, *existing)
	return err
}

func (s *invoiceService) Send(ctx context.Context, id uuid.UUID) error {
	invoice, err := s.invoiceRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	invoice.Status = model.InvoiceStatusSent
	now := time.Now()
	invoice.SentAt = &now
	invoice.UpdatedAt = time.Now()

	_, err = s.invoiceRepository.Update(ctx, *invoice)
	return err
}

func (s *invoiceService) Pay(ctx context.Context, id uuid.UUID) error {
	invoice, err := s.invoiceRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	invoice.Status = model.InvoiceStatusPaid
	now := time.Now()
	invoice.PaidAt = &now
	invoice.UpdatedAt = time.Now()

	_, err = s.invoiceRepository.Update(ctx, *invoice)
	return err
}