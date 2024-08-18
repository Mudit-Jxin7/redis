package storage

import (
	"fmt"
	"strconv"
	"sync"
)

// strings

var store = make(map[string]string)

func Set(key, value string) {
	store[key] = value
}

func Get(key string) string {
	return store[key]
}

func Del(key string) {
	delete(store, key)
}

func Incr(key string) string {
	val, exists := store[key]
	if !exists {
		store[key] = "1"
		return "1"
	}
	intValue, err := strconv.Atoi(val)
	if err != nil {
		return "ERROR: Value is not an integer"
	}
	intValue++
	store[key] = strconv.Itoa(intValue)
	return store[key]
}

func Decr(key string) string {
	val, exists := store[key]
	if !exists {
		store[key] = "1"
		return "1"
	}
	intValue, err := strconv.Atoi(val)
	if err != nil {
		return "ERROR: Value is not an integer"
	}
	intValue--
	store[key] = strconv.Itoa(intValue)
	return store[key]
}

// list

type DB struct {
	sync.RWMutex
	data map[string][]string
}

var db = &DB{
	data: make(map[string][]string),
}

func LPush(key string, values ...string) string {
	db.Lock()
	defer db.Unlock()

	if _, exists := db.data[key]; !exists {
		db.data[key] = []string{}
	}

	db.data[key] = append(values, db.data[key]...)
	return fmt.Sprintf("%d", len(db.data[key]))
}

func RPush(key string, values ...string) string {
	db.Lock()
	defer db.Unlock()

	if _, exists := db.data[key]; !exists {
		db.data[key] = []string{}
	}
	db.data[key] = append(db.data[key], values...)
	return fmt.Sprintf("%d", len(db.data[key]))
}

func LPop(key string) string {
	db.Lock()
	defer db.Unlock()

	if list, exists := db.data[key]; exists && len(list) > 0 {
		val := list[0]
		db.data[key] = list[1:]
		return val
	}
	return "(nil)"
}

func RPop(key string) string {
	db.Lock()
	defer db.Unlock()

	if list, exists := db.data[key]; exists && len(list) > 0 {
		val := list[len(list)-1]
		db.data[key] = list[:len(list)-1]
		return val
	}
	return "(nil)"
}

func LRange(key string, start, stop int) string {
	db.RLock()
	defer db.RUnlock()

	if list, exists := db.data[key]; exists {

		if start < 0 {
			start = 0
		}
		if stop >= len(list) {
			stop = len(list) - 1
		}
		if start > stop {
			return "[]"
		}
		return fmt.Sprintf("%v", list[start:stop+1])
	}
	return "[]"
}

func LLen(key string) string {
	db.RLock()
	defer db.RUnlock()

	if list, exists := db.data[key]; exists {
		return fmt.Sprintf("%d", len(list))
	}
	return "0"
}
