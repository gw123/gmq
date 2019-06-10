package interfaces

type Cache interface {
	Set(key string, data interface{}) error
	Get(key string)
}
