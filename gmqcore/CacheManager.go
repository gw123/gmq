package gmqcore

import (
	"encoding/json"
	"errors"
	"github.com/go-redis/redis"
	"github.com/gw123/gmq"
	"sync"
	"time"
)

type CacheManager struct {
	app         gmq.App
	redisClient *redis.Client
	cacheMap    map[gmq.CacheKey]gmq.CacheRule
	mutex       sync.RWMutex
}

func NewCacheManager(app gmq.App) *CacheManager {
	redisClient, err := app.GetDefaultRedis()
	if err != nil {
		app.Error("CacheManager", "app.GetDefaultRedis() is nil")
	}
	return &CacheManager{
		app:         app,
		redisClient: redisClient,
		cacheMap:    make(map[gmq.CacheKey]gmq.CacheRule),
	}
}

func (c CacheManager) UpdateCache(patten gmq.CacheKey, arg ...interface{}) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	key := gmq.MakeCacheKey(patten, arg...)

	rule, ok := c.cacheMap[patten]
	if !ok {
		return errors.New("not set CacheRule")
	}

	data, err := rule.GetCallback()(arg...)
	if err != nil {
		return err
	}
	tmp, err := json.Marshal(data)
	if err != nil {
		return err
	}
	//if set rule  redisclient , use it first
	redisClient := c.redisClient

	c.app.Debug("CacheManager", "key : %s , tmp  %s", key, tmp)
	return redisClient.Set(key, tmp, time.Hour*24*30).Err()
}

func (c CacheManager) GetCache(out interface{},patten gmq.CacheKey,  arg ...interface{}, ) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	key := gmq.MakeCacheKey(patten, arg...)
	_, ok := c.cacheMap[patten]
	if !ok {
		return errors.New("not set CacheRule")
	}

	redisClient := c.redisClient
	cmd := redisClient.Get(key)
	if cmd.Err() != nil {
		//更新缓存
		if cmd.Err().Error() == "redis: nil" {
			c.app.Pub(gmq.NewUpdateCacheMsg(patten, arg))
		}
		return cmd.Err()
	}
	data := cmd.Val()
	return json.Unmarshal([]byte(data), out)
}

func (c CacheManager) AddCacheRule(rule gmq.CacheRule) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cacheMap[rule.GetCacheKey()] = rule
}

func (c CacheManager) DelCacheRule(patten gmq.CacheKey) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.cacheMap, patten)
}
