package util

import (
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type RedisWrapper struct {
	pool *redis.Pool
}

func NewRedisWrapper(pool *redis.Pool) *RedisWrapper {
	return &RedisWrapper{pool: pool}
}

func (rw *RedisWrapper) Pipeline(ctx context.Context, cb func(conn redis.Conn) error) error {
	conn := rw.pool.Get()
	defer conn.Close()

	err := cb(conn)
	if err != nil {
		return fmt.Errorf("errors.cache Pipeline - detail: %s", err)
	}

	return nil
}

func (rw *RedisWrapper) Exists(ctx context.Context, key string) (bool, error) {
	conn := rw.pool.Get()
	defer conn.Close()

	resp, err := redis.Int(conn.Do("EXISTS", key))
	if err != nil {
		if err == redis.ErrNil {
			return false, nil
		}
		return false, fmt.Errorf("error.redis exist redis - key: %s - detail: %s ", key, err)
	}

	return resp == 1, nil
}

func (rw *RedisWrapper) Del(ctx context.Context, key string) error {
	conn := rw.pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	if err != nil {
		return fmt.Errorf("error.redis del redis - key: %s - detail: %s ", key, err)
	}

	return nil
}
