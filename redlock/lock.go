package redlock

import (
	"context"
	"errors"
	"time"

	"github.com/bsm/redislock"
)

type Lock struct {
	lock *redislock.Lock
}

func (l *Lock) Release(ctx context.Context) error {
	if err := l.lock.Release(ctx); err != nil {
		if errors.Is(err, redislock.ErrLockNotHeld) {
			return nil
		}
		return err
	}
	return nil
}

func (l *Lock) Refresh(ctx context.Context, duration time.Duration) (bool, error) {
	if err := l.lock.Refresh(ctx, duration, nil); err != nil {
		if errors.Is(err, redislock.ErrNotObtained) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
