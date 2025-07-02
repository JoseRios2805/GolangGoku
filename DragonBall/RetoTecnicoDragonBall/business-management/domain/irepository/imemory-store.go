package irepository

import "time"

type IMemoryStore interface {
	Save(key string, expired time.Duration, value interface{}) error
	Get(key string) ([]byte, error)
	Exists(key string) bool
	FlushAll() error
	Delete(key string) error
	GetAll() ([]string, error)
	GetExpiryByName(key string) (int, error)
	GetKeyValue(key string) (string, error)
	IncrementBy(key string, expired time.Duration, count int64) error
	GetKeyValueAndTTLValue(key string) (string, int64, error)
}
