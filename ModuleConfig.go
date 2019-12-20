package gmq2

type ModuleConfig interface {
	GetPath() string
	IsInnerModule() bool
	IsEnable() bool
	GetModuleName() string
	GetItem(string) (value interface{})
	GetItems() (value map[string]interface{})
	GetStringItem(key string) (value string)
	GetItemOrDefault(key, defaultVal string) (value string)
	GetGlobalItem(key string) (value string)
	GetGlobalItems() (value map[string]interface{})
	GetIntItem(key string) int
	GetBoolItem(key string) bool
	SetGlobalConfig(config AppConfig)
	MergeNewConfig(newCofig ModuleConfig) bool
	GetModuleType() string
	GetArrayItem(key string) (value []string)
	GetMapItem(key string) (value map[string]interface{})
	SetItem(key string, value interface{})
}
