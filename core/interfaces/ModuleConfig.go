package interfaces

type ModuleConfig interface {
	GetPath() string
	IsInnerModule() bool
	IsEnable() bool
	GetModuleName() string
	GetItem(key string) (value string)
	GetItemOrDefault(key, defaultVal string) (value string)
	GetGlobalItem(key string) (value string)
	GetItems() (value map[string]interface{})
	GetGlobalItems() (value map[string]interface{})
	SetItem(key string, value interface{})
	GetIntItem(key string) int
	GetBoolItem(key string) bool
	SetGlobalConfig(config AppConfig)
	MergeNewConfig(newCofig ModuleConfig) bool
	GetModuleType() string
	GetArrayItem(key string) (value []string)
	GetMapItem(key string) (value map[string]string)
}
