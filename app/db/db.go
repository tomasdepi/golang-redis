package db

import (
	"sync"
)

type RedisValue struct {
	Val       any
	ExpiresAt int64
	Expires   bool
	Type      int
}

type RedisDB struct {
	data sync.Map
}

func (db *RedisDB) Store(key string, val RedisValue) {
	db.data.Store(key, val)
}

func (db *RedisDB) Load(key string) (RedisValue, bool) {
	v, ok := db.data.Load(key)
	if !ok {
		return RedisValue{}, false
	}
	return v.(RedisValue), true
}

func (db *RedisDB) Delete(key string) {
	db.data.Delete(key)
}
