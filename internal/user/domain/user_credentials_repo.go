package domain

import (
	common_repo "chat/internal/common/repository"
	"context"
)

type UserCredentialsRepo interface {
	GetUserCredentialsByUsername(ctx context.Context, username string) (*UserCredentials, error)
	CreateUserCredentials(ctx context.Context, userCredentials UserCredentials, tx common_repo.Tx) (*UserCredentials, error)
}
