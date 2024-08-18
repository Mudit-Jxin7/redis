package storage

import (
	"strconv"
)

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

func Dec(key string) string {
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
