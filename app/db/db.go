package db

import (
	"fmt"
	"log"
	"sync"
)

const (
	SINGLE_VALUE = 1
	ARRAY_VALUE  = 2
)

type RedisValue struct {
	Val       any
	ExpiresAt int64
	Expires   bool
	Type      int
}

type RedisDB struct {
	data    sync.Map
	waiters sync.Map
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

func (db *RedisDB) AddWaiter(key string, ch chan string) {
	if channels, ok := db.waiters.Load(key); ok {
		new := append(channels.([]chan string), ch)
		log.Printf("I have %s waiter now \n", string(len(new)))
		db.waiters.Store(key, new)
	} else {
		db.waiters.Store(key, []chan string{ch})
	}
}

func (db *RedisDB) PopWaiter(key string) (chan string, error) {
	if channels, ok := db.waiters.Load(key); ok {
		popedCh := channels.([]chan string)[0]
		db.waiters.Store(key, channels.([]chan string)[1:])

		// if len(channels.([]chan string)[1:]) == 0 {
		// 	db.DeleteWaiterEntry(key)
		// }

		return popedCh, nil
	}

	return nil, fmt.Errorf("waiter not found")
}

func (db *RedisDB) DeleteWaiterEntry(key string) {
	db.waiters.Delete(key)
}
