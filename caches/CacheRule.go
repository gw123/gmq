package caches

import (
	"github.com/go-redis/redis"
	"github.com/gw123/GMQ/core/interfaces"
)

type CacheRule struct {
	KeyPatten interfaces.CacheKey
	Callback  func(arg ...interface{}) (interface{}, error)
	Client    *redis.Client
}

func NewCacheRule(keyPatten interfaces.CacheKey, callback func(arg ...interface{}) (interface{}, error), client *redis.Client) *CacheRule {
	return &CacheRule{KeyPatten: keyPatten, Callback: callback, Client: client}
}

func (c CacheRule) GetCacheKey() interfaces.CacheKey {
	return c.KeyPatten
}

func (c CacheRule) GetCallback() func(arg ...interface{}) (interface{}, error) {
	return c.Callback
}

func (c CacheRule) GetRedisClient() *redis.Client {
	return c.Client
}
