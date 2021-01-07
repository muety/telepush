package store

type Store interface {
	Get(key string) interface{}
	GetItems() map[string]interface{}
	Put(key string, value interface{})
	Delete(key string)
	Flush() error
}
