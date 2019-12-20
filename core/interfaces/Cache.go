package interfaces

import (
	"fmt"
	"github.com/go-redis/redis"
)

type CacheKey string

func MakeCacheKey(patten CacheKey, args ...interface{}) string {
	return fmt.Sprintf(string(patten), args...)
}

/***
    KeyPatten  example : group:102 , chapter:102, resource:102
	Callback   call to update  cache  with MakeCacheKey(KeyPatten  , arg...)
	client  redis client to update cache,if nil use app.GetDefaultCache()
*/

type CacheRule interface {
	GetCacheKey() CacheKey
	GetCallback() func(arg ...interface{}) (interface{}, error)
	GetRedisClient() *redis.Client
}

type CacheManager interface {
	UpdateCache(patten CacheKey, arg ...interface{}) error
	GetCache(out interface{}, patten CacheKey, arg ...interface{}) error
	AddCacheRule(rule CacheRule)
	DelCacheRule(patten CacheKey)
}
