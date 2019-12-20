package gmq

import (
	"fmt"
)

type CacheKey string
type UpdateCacheCallback func(arg ...interface{}) (interface{}, error)
func MakeCacheKey(patten CacheKey, args ...interface{}) string {
	return fmt.Sprintf(string(patten), args...)
}

/***
    KeyPatten  example : group:102 , chapter:102, resource:102
	Callback   call to update  cache  with MakeCacheKey(KeyPatten  , arg...)
*/

type CacheRule interface {
	GetCacheKey() CacheKey
	GetCallback() UpdateCacheCallback
}

type CacheManager interface {
	UpdateCache(patten CacheKey, arg ...interface{}) error
	GetCache(out interface{}, patten CacheKey, arg ...interface{}) error
	AddCacheRule(rule CacheRule)
	DelCacheRule(patten CacheKey)
}
