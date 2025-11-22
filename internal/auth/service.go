package auth

import (
	"errors"
	"time"

	"github.com/adityanuriskandar17/HRIS-BE/internal/domain/model"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")
	ErrRefreshTokenRevoked  = errors.New("refresh token revoked")
)

type Service struct {
	db         *gorm.DB
	jwtSecret  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewService(db *gorm.DB, secret string, accessTTL, refreshTTL time.Duration) *Service {
	return &Service{
		db:         db,
		jwtSecret:  secret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *Service) IssueTokenPair(user *model.UserAccount) (*TokenPair, error) {
	if user == nil {
		return nil, errors.New("user required")
	}

	accessToken, accessExp, err := SignAccessToken(user.ID.String(), string(user.Role), s.jwtSecret, s.accessTTL, user.TenantID.String())
	if err != nil {
		return nil, err
	}

	refreshPlain, refreshHash, refreshExp, err := GenerateRefreshToken(s.refreshTTL)
	if err != nil {
		return nil, err
	}

	rt := &model.RefreshToken{
		UserID:    user.ID,
		TokenHash: refreshHash,
		ExpiresAt: refreshExp,
	}
	if err := s.db.Create(rt).Error; err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshPlain,
		AccessExpiresAt:  accessExp,
		RefreshExpiresAt: refreshExp,
	}, nil
}

func (s *Service) RotateRefreshToken(plain string) (*TokenPair, *model.UserAccount, error) {
	if plain == "" {
		return nil, nil, ErrRefreshTokenNotFound
	}

	hash := HashRefreshToken(plain)
	var (
		pair TokenPair
		user model.UserAccount
	)

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var stored model.RefreshToken
		if err := tx.Preload("User").Where("token_hash = ?", hash).First(&stored).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrRefreshTokenNotFound
			}
			return err
		}

		now := time.Now().UTC()
		if stored.RevokedAt != nil {
			return ErrRefreshTokenRevoked
		}
		if now.After(stored.ExpiresAt) {
			return ErrRefreshTokenExpired
		}

		if err := tx.Model(&stored).Update("revoked_at", now).Error; err != nil {
			return err
		}

		user = stored.User

		accessToken, accessExp, err := SignAccessToken(user.ID.String(), string(user.Role), s.jwtSecret, s.accessTTL, user.TenantID.String())
		if err != nil {
			return err
		}

		refreshPlain, refreshHash, refreshExp, err := GenerateRefreshToken(s.refreshTTL)
		if err != nil {
			return err
		}

		newRT := &model.RefreshToken{
			UserID:    user.ID,
			TokenHash: refreshHash,
			ExpiresAt: refreshExp,
		}

		if err := tx.Create(newRT).Error; err != nil {
			return err
		}

		pair = TokenPair{
			AccessToken:      accessToken,
			RefreshToken:     refreshPlain,
			AccessExpiresAt:  accessExp,
			RefreshExpiresAt: refreshExp,
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return &pair, &user, nil
}
