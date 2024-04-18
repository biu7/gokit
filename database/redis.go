package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

func NewRedis(link string, trace bool) (*redis.Client, error) {
	opts, err := redis.ParseURL(link)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis url: %w", err)
	}
	rdb := redis.NewClient(opts)
	if trace {
		// 链路追踪
		if err = redisotel.InstrumentTracing(rdb); err != nil {
			return nil, fmt.Errorf("failed to instrument redis tracing: %w", err)
		}
	}

	err = rdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}
	return rdb, nil
}

func RedisNotFound(err error) bool {
	return errors.Is(err, redis.Nil)
}
