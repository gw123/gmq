package gmq2

import (
	"github.com/go-redis/redis"
)

/**
const (
Resource         interfaces.CacheKey = "Resource:"
Group            interfaces.CacheKey = "Group:"
GroupLatestNews  interfaces.CacheKey = "GroupLatestNews:%d"
Chapter          interfaces.CacheKey = "Chapter:"
GroupTag         interfaces.CacheKey = "GroupTag:"
)
*/

type RedisCacheRule struct {
	KeyPatten CacheKey
	Callback  UpdateCacheCallback
	Client    *redis.Client
}

func NewCacheRule(keyPatten CacheKey, callback UpdateCacheCallback, client *redis.Client) *RedisCacheRule {
	return &RedisCacheRule{KeyPatten: keyPatten, Callback: callback, Client: client}
}

func (c RedisCacheRule) GetCacheKey() CacheKey {
	return c.KeyPatten
}

func (c RedisCacheRule) GetCallback() func(arg ...interface{}) (interface{}, error) {
	return c.Callback
}

func (c RedisCacheRule) GetRedisClient() *redis.Client {
	return c.Client
}
