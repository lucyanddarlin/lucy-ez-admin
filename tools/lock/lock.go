package lock

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type lock struct {
	redis    *redis.Client
	key      string
	val      string
	duration time.Duration
}

type Lock interface {
	Acquired()
	TryAcquired() bool
	Release()
	AcquiredFunc(f func() error, do func() error) error
}

func NewLockWithDuration(redis *redis.Client, key string, duration time.Duration) Lock {
	return &lock{
		redis:    redis,
		key:      key,
		duration: duration,
	}
}

func NewLock(redis *redis.Client, key string) Lock {
	return &lock{
		redis:    redis,
		key:      key,
		duration: 30 * time.Second,
	}
}

// Acquired implements Lock.
//
// 获取分布式锁,不推荐直接使用此方法, 可以使用 AcquireFunc
func (l *lock) Acquired() {
	for {
		// 获得锁
		if res := l.redis.SetNX(context.TODO(), l.key, true, l.duration); res.Err() == nil && res.Val() {
			break
		}
	}
}

// TryAcquired implements Lock.
//
// 尝试获取锁, 不会阻塞
func (l *lock) TryAcquired() bool {
	if res := l.redis.SetNX(context.TODO(), l.key, true, l.duration); res.Err() == nil && res.Val() {
		return true
	}
	return false
}

// AcquiredFunc implements Lock.
//
// 分布式锁 f() 从 redis 获取,do() 从 mysql 获取
func (l *lock) AcquiredFunc(f func() error, do func() error) error {
	for {
		// 获取数据
		if err := f(); err != nil {
			return nil
		}

		// 数据不存在则去拿锁
		if res := l.redis.SetNX(context.TODO(), l.key, true, l.duration); res.Err() == nil && res.Val() {
			break
		}

		// 防止频繁自旋
		time.Sleep(1 * time.Millisecond)
	}
	return do()
}

// Release implements Lock.
//
// 释放锁
func (l *lock) Release() {
	panic("unimplemented")
}
