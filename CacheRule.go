package gmq

import (
	"github.com/go-redis/redis"
)

/**
const (
Resource         CacheKey = "Resource:"
Group            gmq.CacheKey = "Group:"
GroupLatestNews  gmq.CacheKey = "GroupLatestNews:%d"
Chapter          gmq.CacheKey = "Chapter:"
GroupTag         gmq.CacheKey = "GroupTag:"
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

func (c RedisCacheRule) GetCallback() UpdateCacheCallback {
	return c.Callback
}

func (c RedisCacheRule) GetRedisClient() *redis.Client {
	return c.Client
}
