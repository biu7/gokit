package redlock

import (
	"context"
	"errors"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

type Locker struct {
	client *redislock.Client
}

func NewLocker(rds *redis.Client) *Locker {
	client := redislock.New(rds)
	return &Locker{
		client: client,
	}
}

func (l *Locker) TryLock(ctx context.Context, resource string, duration time.Duration) (*Lock, bool, error) {
	lock, err := l.client.Obtain(ctx, "locker:"+resource, duration, nil)
	if err != nil {
		if errors.Is(err, redislock.ErrNotObtained) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &Lock{lock: lock}, true, nil
}
