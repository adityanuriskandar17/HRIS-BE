package repository

import (
	"context"
	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	"github.com/google/uuid"
)

// SubscriptionRepository defines the interface for subscription data operations
type SubscriptionRepository interface {
	FindAll(ctx context.Context) ([]model.Subscription, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error)
	FindByTenantID(ctx context.Context, tenantID uuid.UUID) ([]model.Subscription, error)
	Create(ctx context.Context, subscription model.Subscription) (*model.Subscription, error)
	Update(ctx context.Context, subscription model.Subscription) (*model.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
