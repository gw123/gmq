package redisModule

import "github.com/gw123/GMQ/core/interfaces"

type RedisModuleProvider struct {
	module interfaces.Module
}

func NewRedisModuleProvider() *RedisModuleProvider {
	this := new(RedisModuleProvider)
	return this
}

func (this *RedisModuleProvider) GetModuleName() string {
	return "RedisModule"
}

func (this *RedisModuleProvider) Register() {
}

func (this *RedisModuleProvider) GetModule() interfaces.Module {
	if this.module != nil {
		return this.module
	}
	this.module = NewRedisModule()
	return this.module
}

func (this *RedisModuleProvider) GetNewModule() interfaces.Module {
	this.module = NewRedisModule()
	return this.module
}
