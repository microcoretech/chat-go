package domain

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"golang.org/x/crypto/bcrypt"

	"mbobrovskyi/chat-go/internal/common/domain"
	"mbobrovskyi/chat-go/internal/common/errors"
	"mbobrovskyi/chat-go/internal/common/repository"
	usererrors "mbobrovskyi/chat-go/internal/user/errors"
)

type AuthServiceImpl struct {
	baseRepo            repository.BaseRepo
	userRepo            UserRepo
	userCredentialsRepo UserCredentialsRepo
	sessionRepository   SessionRepository
}

func (s *AuthServiceImpl) createToken(ctx context.Context, userCredentials UserCredentials) (*Token, error) {
	const (
		tokenLength = 32
		tokenTTL    = 30 * 24 * time.Hour // 30 days
	)

	token := &Token{
		Token:   s.generateSecureToken(tokenLength),
		ExpIn:   tokenTTL,
		ExpTime: time.Now().Add(tokenTTL),
		User:    userCredentials.User,
	}

	return token, nil
}

func (s *AuthServiceImpl) createTokenAndSetSession(ctx context.Context, userCredentials UserCredentials) (*Token, error) {
	token, err := s.createToken(ctx, userCredentials)
	if err != nil {
		return nil, err
	}

	if err := s.sessionRepository.SetSession(
		ctx, token.Token, domain.UserToSession(userCredentials.User), token.ExpIn); err != nil {
		return nil, err
	}

	return token, err
}

func (s *AuthServiceImpl) SignIn(ctx context.Context, username, password string) (*Token, error) {
	userCredentials, err := s.userCredentialsRepo.GetUserCredentialsByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if userCredentials == nil || userCredentials.User.ID == 0 {
		return nil, errors.NewUnauthorizedError("invalid credentials")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userCredentials.Password), []byte(password)); err != nil {
		return nil, errors.NewUnauthorizedError("invalid credentials")
	}

	token, err := s.createTokenAndSetSession(ctx, *userCredentials)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *AuthServiceImpl) SignUp(ctx context.Context, newUser domain.User, password string) (*Token, error) {
	userAlreadyCreatedErrorData := make(map[string]any)

	userCredentialsByEmail, err := s.userCredentialsRepo.GetUserCredentialsByUsername(ctx, newUser.Email)
	if err != nil {
		return nil, err
	}
	if userCredentialsByEmail != nil {
		userAlreadyCreatedErrorData["email"] = newUser.Email
	}

	userCredentialsByUsername, err := s.userCredentialsRepo.GetUserCredentialsByUsername(ctx, newUser.Username)
	if err != nil {
		return nil, err
	}
	if userCredentialsByUsername != nil {
		userAlreadyCreatedErrorData["username"] = newUser.Username
	}

	if len(userAlreadyCreatedErrorData) > 0 {
		return nil, usererrors.NewUserAlreadyCreatedError(userAlreadyCreatedErrorData)
	}

	tx, err := s.baseRepo.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	createdUser, err := s.userRepo.CreateUser(ctx, newUser, tx)
	if err != nil {
		return nil, err
	}
	if createdUser == nil {
		return nil, usererrors.NewUserNotCreatedError()
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUserCredentials := UserCredentials{UserID: createdUser.ID, Password: string(hash)}
	createdUserCredentials, err := s.userCredentialsRepo.CreateUserCredentials(ctx, newUserCredentials, tx)
	if err != nil {
		return nil, err
	}
	if createdUserCredentials == nil {
		return nil, usererrors.NewUserNotCreatedError()
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	token, err := s.createTokenAndSetSession(ctx, *createdUserCredentials)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *AuthServiceImpl) SignOut(ctx context.Context, userID uint64, token string) error {
	if err := s.sessionRepository.DeleteSession(ctx, userID, token); err != nil {
		return err
	}
	return nil
}

func (s *AuthServiceImpl) GetSession(ctx context.Context, token string) (*domain.Session, error) {
	session, err := s.sessionRepository.GetSession(ctx, token)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *AuthServiceImpl) generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func NewAuthServiceImpl(
	baseRepo repository.BaseRepo,
	userRepo UserRepo,
	userCredentialsRepo UserCredentialsRepo,
	sessionRepository SessionRepository,
) *AuthServiceImpl {
	return &AuthServiceImpl{
		baseRepo:            baseRepo,
		userRepo:            userRepo,
		userCredentialsRepo: userCredentialsRepo,
		sessionRepository:   sessionRepository,
	}
}
