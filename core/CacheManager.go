package core

import (
	"encoding/json"
	"errors"
	"github.com/go-redis/redis"
	"github.com/gw123/GMQ/common/gmsg"
	"github.com/gw123/GMQ/core/interfaces"
	"sync"
	"time"
)

type CacheManager struct {
	app         interfaces.App
	redisClient *redis.Client
	cacheMap    map[interfaces.CacheKey]interfaces.CacheRule
	mutex       sync.RWMutex
}

func NewCacheManager(app interfaces.App) *CacheManager {
	redisClient, err := app.GetDefaultRedis()
	if err != nil {
		app.Error("CacheManager", "app.GetDefaultRedis() is nil")
	}
	return &CacheManager{
		app:         app,
		redisClient: redisClient,
		cacheMap:    make(map[interfaces.CacheKey]interfaces.CacheRule),
	}
}

func (c CacheManager) UpdateCache(patten interfaces.CacheKey, arg ...interface{}) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	key := interfaces.MakeCacheKey(patten, arg...)

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
	if rule.GetRedisClient() != nil {
		redisClient = rule.GetRedisClient()
	}
	c.app.Debug("CacheManager", "key : %s , tmp  %s", key, tmp)
	return redisClient.Set(key, tmp, time.Hour*24*30).Err()
}

func (c CacheManager) GetCache(patten interfaces.CacheKey, out interface{}, arg ...interface{}, ) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	key := interfaces.MakeCacheKey(patten, arg...)
	rule, ok := c.cacheMap[patten]
	if !ok {
		return errors.New("not set CacheRule")
	}

	//if set rule  redisclient , use it first
	redisClient := c.redisClient
	if rule.GetRedisClient() != nil {
		redisClient = rule.GetRedisClient()
	}

	cmd := redisClient.Get(key)
	if cmd.Err() != nil {
		//更新缓存
		if cmd.Err().Error() == "redis: nil" {
			c.app.Pub(gmsg.NewUpdateCacheMsg(patten, arg))
		}
		return cmd.Err()
	}
	data := cmd.Val()
	return json.Unmarshal([]byte(data), out)
}

func (c CacheManager) AddCacheRule(rule interfaces.CacheRule) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cacheMap[rule.GetCacheKey()] = rule
}

func (c CacheManager) DelCacheRule(patten interfaces.CacheKey) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.cacheMap, patten)
}
