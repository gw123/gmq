package interfaces

type ModuleConfig interface {
	GetPath() string
	IsInnerModule() bool
	IsEnable() bool
	GetModuleName() string
	GetItem(key string) (value string)
	GetGlobalItem(key string) (value string)
	GetItems() (value map[string]string)
	GetGlobalItems() (value map[string]string)
	SetItem(key, value string)
}
