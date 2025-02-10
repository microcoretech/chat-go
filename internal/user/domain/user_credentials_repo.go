package domain

import (
	"context"

	common_repo "mbobrovskyi/chat-go/internal/common/repository"
)

type UserCredentialsRepo interface {
	GetUserCredentialsByUsername(ctx context.Context, username string) (*UserCredentials, error)
	CreateUserCredentials(ctx context.Context, userCredentials UserCredentials, tx common_repo.Tx) (*UserCredentials, error)
}
