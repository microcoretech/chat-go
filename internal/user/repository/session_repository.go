package repository

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"mbobrovskyi/chat-go/internal/common/domain"
	commonerrors "mbobrovskyi/chat-go/internal/common/errors"
	"mbobrovskyi/chat-go/internal/user/common"
)

type SessionRepository struct {
	redisClient *redis.Client
}

func (r *SessionRepository) GetSession(ctx context.Context, token string) (*domain.Session, error) {
	tokenKey := r.getTokenKey(token)

	sessionJSON, err := r.redisClient.Get(ctx, tokenKey).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, commonerrors.NewDatabaseError(common.UserDomain, err, "error on get session from redis")
	}

	var session domain.Session

	if err = json.Unmarshal(sessionJSON, &session); err != nil {
		return nil, commonerrors.NewDatabaseError(common.UserDomain, err, "error on unmarshal session")
	}

	return &session, nil
}

func (r *SessionRepository) SetSession(ctx context.Context, token string, session domain.Session, expTime time.Duration) error {
	sessionBytes, err := json.Marshal(session)
	if err != nil {
		return commonerrors.NewDatabaseError(common.UserDomain, err, "error on marshal session")
	}

	tokenKey := r.getTokenKey(token)
	userKey := r.getUserKey(session.User.ID)

	if err = r.redisClient.Set(ctx, tokenKey, sessionBytes, expTime).Err(); err != nil {
		return commonerrors.NewDatabaseError(common.UserDomain, err, "error on set session")
	}

	unixNow := time.Now().UTC().Unix()
	unixNowStr := strconv.FormatInt(unixNow, 10)

	if err = r.redisClient.ZRemRangeByScore(ctx, userKey, "-inf", unixNowStr).Err(); err != nil {
		return commonerrors.NewDatabaseError(common.UserDomain, err, "error on remove tokens from tokens list")
	}
	if expTime == redis.KeepTTL {
		return nil
	}

	params := redis.Z{
		Score:  float64(time.Now().UTC().Add(expTime).Unix()),
		Member: token,
	}

	err = r.redisClient.ZAdd(ctx, userKey, params).Err()
	if err != nil {
		return commonerrors.NewDatabaseError(common.UserDomain, err, "error on add token to tokens list")
	}
	return nil
}

func (r *SessionRepository) DeleteSession(ctx context.Context, userID uint64, token string) error {
	tokenKey := r.getTokenKey(token)
	userKey := r.getUserKey(userID)

	if err := r.redisClient.Del(ctx, tokenKey).Err(); err != nil {
		return commonerrors.NewDatabaseError(common.UserDomain, err, "error on delete token")
	}
	if err := r.redisClient.ZRem(ctx, userKey, token).Err(); err != nil {
		return commonerrors.NewDatabaseError(common.UserDomain, err, "error on delete token from token list")
	}
	return nil
}

func (r *SessionRepository) GetTokensByUserID(ctx context.Context, userID uint64) ([]string, error) {
	var tokens []string
	userKey := r.getUserKey(userID)

	params := &redis.ZRangeBy{
		Min: strconv.FormatInt(time.Now().Unix(), 10),
		Max: "+inf",
	}

	if err := r.redisClient.ZRangeByScore(ctx, userKey, params).ScanSlice(&tokens); err != nil {
		return nil, commonerrors.NewDatabaseError(common.UserDomain, err, "error on get tokens from tokens list")
	}

	return tokens, nil
}

func (r *SessionRepository) getTokenKey(token string) string {
	return "token:" + token
}

func (r *SessionRepository) getUserKey(userID uint64) string {
	return "user:" + strconv.FormatUint(userID, 10)
}

func NewSessionRepo(
	redisClient *redis.Client,
) *SessionRepository {
	return &SessionRepository{
		redisClient: redisClient,
	}
}
