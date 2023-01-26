package cache

type Cache interface {
	Get(key string) interface{}
	Set(key string, value interface{})
	Del(key string)
	Flush() error
}
