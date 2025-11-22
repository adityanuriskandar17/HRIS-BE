package authctx

import (
	"context"

	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	"github.com/google/uuid"
)

type contextKey string

const userContextKey contextKey = "auth.user"

type AuthenticatedUser struct {
	ID       uuid.UUID
	Role     model.UserRole
	TenantID string
}

func WithUser(ctx context.Context, user AuthenticatedUser) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func CurrentUser(ctx context.Context) (AuthenticatedUser, bool) {
	val := ctx.Value(userContextKey)
	if val == nil {
		return AuthenticatedUser{}, false
	}
	user, ok := val.(AuthenticatedUser)
	return user, ok
}
