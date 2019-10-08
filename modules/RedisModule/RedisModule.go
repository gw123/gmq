package redisModule

import (
	"github.com/gw123/GMQ/core/interfaces"
	"github.com/gw123/GMQ/modules/base"
)

type RedisModule struct {
	base.BaseModule
}

func NewRedisModule() *RedisModule {
	this := new(RedisModule)
	return this
}

func (this *RedisModule) Init(app interfaces.App, config interfaces.ModuleConfig) error {
	this.BaseModule.Init(app, this, config)
	//app.Sub("debug", this)
	return nil
}

func (this *RedisModule) Handle(event interfaces.Event) error {
	return nil
}

func (this *RedisModule) Watch(index int) {
	return
}
