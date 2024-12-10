package redis

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/lucyanddarlin/lucy-ez-admin/config"
)

type rd struct {
	mu sync.RWMutex
	m  map[string]*redis.Client
}

type Redis interface {
	// Get
	//
	//  @Description: 获取指定名称的 redis 实例, 如果实例不存在则返回报错
	//  @param 实例名称
	//  @return *redis.Client
	//  @return error
	Get(name string) (*redis.Client, error)
	// GetRedis
	//
	//  @Description: 获取指定名称的 redis 实例, 如果实例不存在则返回 nil
	//  @param 实例名称
	//  @return *redis.Client
	//  @return error
	GetRedis(name string) *redis.Client
}

func New(rc []config.Redis) Redis {
	redisIns := rd{
		mu: sync.RWMutex{},
		m:  make(map[string]*redis.Client),
	}

	redisIns.mu.Lock()
	defer redisIns.mu.Unlock()

	for _, conf := range rc {
		if !conf.Enable {
			continue
		}

		client := redis.NewClient(&redis.Options{
			Addr:     conf.Host,
			Username: conf.Username,
			Password: conf.Password,
		})
		if err := client.Ping(context.TODO()).Err(); err != nil {
			panic(fmt.Sprintf("redis 初始化失败: %v", err.Error()))
		}

		redisIns.m[conf.Name] = client
	}

	return &redisIns
}

// Get implements Redis.
func (r *rd) Get(name string) (*redis.Client, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.m[name] == nil {
		return nil, errors.New("no exist redis")
	}
	return r.m[name], nil
}

// GetRedis implements Redis.
func (r *rd) GetRedis(name string) *redis.Client {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.m[name]
}
