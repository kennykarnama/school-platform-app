package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/kennykarnama/school-user-api/domain/models/auth"
	"github.com/kennykarnama/school-user-api/util"
)

const (
	// AccessTokenKey forms
	// access-token:user:[uuid]:
	AccessTokenKey = "token:%s"
	// RefreshTokenKey forms
	// refresh-token:user:[uuid]
	RefreshTokenKey = "refresh-token:%s"
)

type repository struct {
	redisWrapper *util.RedisWrapper
}

func NewRepository(redisWrapper *util.RedisWrapper) *repository {
	return &repository{redisWrapper: redisWrapper}
}

func (r *repository) SaveJWTSession(ctx context.Context, userID string, metadata *auth.JWTMetadata) error {
	now := time.Now().UTC()
	err := r.redisWrapper.Pipeline(ctx, func(conn redis.Conn) error {
		conn.Send("MULTI")

		fmt.Println(metadata.AccessTokenExpiresTime().Sub(now))
		conn.Send("SETEX", r.tokenKey(AccessTokenKey, metadata.TokenID), int64(metadata.AccessTokenExpiresTime().Sub(now).Seconds()), userID)
		conn.Send("SETEX", r.tokenKey(RefreshTokenKey, metadata.TokenID), int64(metadata.RefreshTokenExpiresTime().Sub(now).Seconds()), userID)

		_, err := conn.Do("EXEC")
		return err
	})

	if err != nil {
		return fmt.Errorf("action=repo.saveJWTSession userID=%v err=%v", userID, err)
	}

	return nil
}

func (r *repository) IsAccessTokenExist(ctx context.Context, tokenID string) (bool, error) {
	return r.redisWrapper.Exists(ctx, r.tokenKey(AccessTokenKey, tokenID))
}

func (r *repository) IsRefreshTokenExist(ctx context.Context, tokenID string) (bool, error) {
	return r.redisWrapper.Exists(ctx, r.tokenKey(RefreshTokenKey, tokenID))
}

func (r *repository) DeleteJWTSession(ctx context.Context, token string) error {
	err := r.redisWrapper.Pipeline(ctx, func(conn redis.Conn) error {
		conn.Send("MULTI")

		conn.Send("DEL", r.tokenKey(AccessTokenKey, token))
		conn.Send("DEL", r.tokenKey(RefreshTokenKey, token))

		_, err := conn.Do("EXEC")
		return err
	})

	if err != nil {
		return fmt.Errorf("action=repo.saveJWTSession tokenID=%v err=%v", token, err)
	}

	return nil
}

func (r *repository) tokenKey(form string, tokenID string) string {
	return fmt.Sprintf(form, tokenID)
}
