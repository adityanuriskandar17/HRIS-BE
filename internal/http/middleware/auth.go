package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/adityanuriskandar17/HRIS-BE/internal/auth"
	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	"github.com/adityanuriskandar17/HRIS-BE/internal/http/authctx"
	"github.com/adityanuriskandar17/HRIS-BE/internal/http/envelope"
	"gorm.io/gorm"
)

type Authenticator struct {
	Secret string
	DB     *gorm.DB
}

func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			envelope.Error(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "missing authorization header", nil)
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			envelope.Error(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "invalid authorization header", nil)
			return
		}

		claims, err := auth.ParseAccessToken(strings.TrimSpace(parts[1]), a.Secret)
		if err != nil {
			envelope.Error(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "invalid token", nil)
			return
		}

		uid, err := strconv.ParseUint(claims.Subject, 10, 64)
		if err != nil {
			envelope.Error(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "invalid token subject", nil)
			return
		}

		var user model.UserAccount
		if err := a.DB.First(&user, uid).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				envelope.Error(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "user not found", nil)
				return
			}
			envelope.Error(w, r, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "failed to load user", nil)
			return
		}

		if user.Status == 0 {
			envelope.Error(w, r, http.StatusForbidden, "FORBIDDEN", "account disabled", nil)
			return
		}

		ctx := authctx.WithUser(r.Context(), authctx.AuthenticatedUser{ID: user.ID, Role: user.Role, TenantID: claims.TenantID})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireRoles(roles ...model.UserRole) func(http.Handler) http.Handler {
	allowed := make(map[model.UserRole]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := authctx.CurrentUser(r.Context())
			if !ok {
				envelope.Error(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "unauthorized", nil)
				return
			}
			if _, ok := allowed[user.Role]; !ok {
				envelope.Error(w, r, http.StatusForbidden, "FORBIDDEN", "forbidden", nil)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
