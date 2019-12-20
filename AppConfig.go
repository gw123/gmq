package gmq2

type AppConfig interface {
	GetItem(key string) (value string)
	SetItem(key, value string)
	GetItems() (value map[string]interface{})
}
