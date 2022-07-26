package cache

import "time"

type Cache interface {
	Get(key string) (string, error)
	Set(key string, val interface{}, timeout time.Duration) error
	IsExist(key string) bool
	Delete(key string) error
}
