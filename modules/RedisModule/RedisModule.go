package redisModule

import (
	"errors"
	"github.com/gw123/GMQ/common/gmsg"
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/modules/base"
)

type RedisModule struct {
	base.BaseModule
	cacheManager interfaces.CacheManager
}

func NewRedisModule() *RedisModule {
	this := new(RedisModule)
	return this
}

func (m *RedisModule) Init(app interfaces.App, config interfaces.ModuleConfig) error {
	if err := m.BaseModule.Init(app, m, config); err != nil {
		return err
	}
	var err error
	m.cacheManager, err = app.GetCacheManager()
	if err != nil {
		return err
	}
	app.Sub("updateCache", m)
	return nil
}

func (m *RedisModule) Handle(msg interfaces.Msg) error {
	switch msg.GetEventName() {
	case "updateCache":
		updateCacheMsg, ok := msg.(*gmsg.UpdateCacheMsg)
		if !ok {
			return errors.New("消息类型错误")
		}
		m.Info("update Cache ,key: %s ", interfaces.MakeCacheKey(updateCacheMsg.Cachekey, updateCacheMsg.Arguments...))
		return m.cacheManager.UpdateCache(updateCacheMsg.Cachekey, updateCacheMsg.Arguments...)
	}
	return nil
}

func (m *RedisModule) Watch(index int) {
	return
}
