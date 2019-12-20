package gmq

type ModuleProvider interface {
	//获取模块名
	GetModuleName() string

	//注册模块
	Register()

	//获取模块实例
	GetModule() Module

	//强制获取新的模块
	GetNewModule() Module
}
