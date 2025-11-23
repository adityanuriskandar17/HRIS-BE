package services

import (
	"context"
	"time"

	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/repository"
	"github.com/google/uuid"
)

type subscriptionService struct {
	subscriptionRepository repository.SubscriptionRepository
	planRepository         repository.PlanRepository
}

func NewSubscriptionService(subscriptionRepository repository.SubscriptionRepository, planRepository repository.PlanRepository) SubscriptionService {
	return &subscriptionService{
		subscriptionRepository: subscriptionRepository,
		planRepository:         planRepository,
	}
}

type SubscriptionService interface {
	GetAll(ctx context.Context) ([]*model.Subscription, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error)
	GetByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*model.Subscription, error)
	Create(ctx context.Context, subscription *model.Subscription) error
	Update(ctx context.Context, id uuid.UUID, subscription *model.Subscription) error
	Cancel(ctx context.Context, id uuid.UUID) error
	Renew(ctx context.Context, id uuid.UUID, months int) error
}

func (s *subscriptionService) GetAll(ctx context.Context) ([]*model.Subscription, error) {
	subscriptions, err := s.subscriptionRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Subscription, len(subscriptions))
	for i, sub := range subscriptions {
		result[i] = &sub
	}

	return result, nil
}

func (s *subscriptionService) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	return s.subscriptionRepository.FindByID(ctx, id)
}

func (s *subscriptionService) GetByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*model.Subscription, error) {
	subscriptions, err := s.subscriptionRepository.FindByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	result := make([]*model.Subscription, len(subscriptions))
	for i, sub := range subscriptions {
		result[i] = &sub
	}

	return result, nil
}

func (s *subscriptionService) Create(ctx context.Context, subscription *model.Subscription) error {
	// Set default values
	subscription.ID = uuid.New()
	subscription.Status = model.SubscriptionStatusActive
	subscription.CreatedAt = time.Now()
	subscription.UpdatedAt = time.Now()

	// Calculate end date based on plan billing cycle
	if subscription.PlanID != uuid.Nil {
		plan, err := s.planRepository.FindByID(ctx, subscription.PlanID)
		if err != nil {
			return err
		}
		// Use BillingCycle to determine duration
		months := 1 // Default to 1 month
		if plan.BillingCycle == "yearly" {
			months = 12
		}
		subscription.EndDate = subscription.StartDate.AddDate(0, months, 0)
	}

	_, err := s.subscriptionRepository.Create(ctx, *subscription)
	return err
}

func (s *subscriptionService) Update(ctx context.Context, id uuid.UUID, subscription *model.Subscription) error {
	// Get existing subscription
	existing, err := s.subscriptionRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Update fields
	existing.PlanID = subscription.PlanID
	existing.StartDate = subscription.StartDate
	existing.EndDate = subscription.EndDate
	existing.Status = subscription.Status
	existing.UpdatedAt = time.Now()

	// Recalculate end date if plan changed
	if subscription.PlanID != uuid.Nil && subscription.PlanID != existing.PlanID {
		plan, err := s.planRepository.FindByID(ctx, subscription.PlanID)
		if err != nil {
			return err
		}
		// Use BillingCycle to determine duration
		months := 1 // Default to 1 month
		if plan.BillingCycle == "yearly" {
			months = 12
		}
		existing.EndDate = existing.StartDate.AddDate(0, months, 0)
	}

	_, err = s.subscriptionRepository.Update(ctx, *existing)
	return err
}

func (s *subscriptionService) Cancel(ctx context.Context, id uuid.UUID) error {
	subscription, err := s.subscriptionRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	subscription.Status = model.SubscriptionStatusCanceled
	subscription.UpdatedAt = time.Now()

	_, err = s.subscriptionRepository.Update(ctx, *subscription)
	return err
}

func (s *subscriptionService) Renew(ctx context.Context, id uuid.UUID, months int) error {
	subscription, err := s.subscriptionRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Extend the end date
	subscription.EndDate = subscription.EndDate.AddDate(0, months, 0)
	subscription.Status = model.SubscriptionStatusActive
	subscription.UpdatedAt = time.Now()

	_, err = s.subscriptionRepository.Update(ctx, *subscription)
	return err
}
