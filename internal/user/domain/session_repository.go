package domain

import (
	"context"
	"time"

	"mbobrovskyi/chat-go/internal/common/domain"
)

type SessionRepository interface {
	GetSession(ctx context.Context, token string) (*domain.Session, error)
	SetSession(ctx context.Context, token string, session domain.Session, expTime time.Duration) error
	DeleteSession(ctx context.Context, userID uint64, token string) error
	GetTokensByUserID(ctx context.Context, userID uint64) ([]string, error)
}
