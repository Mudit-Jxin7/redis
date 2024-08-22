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
	hash map[string]map[string]string
	set  map[string]map[string]struct{}
}

var db = &DB{
	data: make(map[string][]string),
	hash: make(map[string]map[string]string),
	set:  make(map[string]map[string]struct{}),
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

// hash

func HSet(key, field, value string) string {
	db.Lock()
	defer db.Unlock()

	if _, exists := db.hash[key]; !exists {
		db.hash[key] = make(map[string]string)
	}
	db.hash[key][field] = value
	return "OK"
}

func HGet(key, field string) string {
	db.RLock()
	defer db.RUnlock()

	if fields, exists := db.hash[key]; exists {
		if val, ok := fields[field]; ok {
			return val
		}
	}
	return "(nil)"
}

func HMSet(key string, fieldValues map[string]string) string {
	db.Lock()
	defer db.Unlock()

	if _, exists := db.hash[key]; !exists {
		db.hash[key] = make(map[string]string)
	}
	for field, value := range fieldValues {
		db.hash[key][field] = value
	}
	return "OK"
}

func HMGet(key string, fields ...string) []string {
	db.RLock()
	defer db.RUnlock()

	var result []string
	if fieldMap, exists := db.hash[key]; exists {
		for _, field := range fields {
			if val, ok := fieldMap[field]; ok {
				result = append(result, val)
			} else {
				result = append(result, "(nil)")
			}
		}
	} else {
		for range fields {
			result = append(result, "(nil)")
		}
	}
	return result
}

func HGetAll(key string) map[string]string {
	db.RLock()
	defer db.RUnlock()

	if fields, exists := db.hash[key]; exists {
		return fields
	}
	return map[string]string{}
}

func HDel(key string, fields ...string) int {
	db.Lock()
	defer db.Unlock()

	count := 0
	if fieldMap, exists := db.hash[key]; exists {
		for _, field := range fields {
			if _, ok := fieldMap[field]; ok {
				delete(fieldMap, field)
				count++
			}
		}

		if len(fieldMap) == 0 {
			delete(db.hash, key)
		}
	}
	return count
}

// sets

func SAdd(key string, members ...string) int {
	db.Lock()
	defer db.Unlock()

	if _, exists := db.set[key]; !exists {
		db.set[key] = make(map[string]struct{})
	}

	added := 0
	for _, member := range members {
		if _, exists := db.set[key][member]; !exists {
			db.set[key][member] = struct{}{}
			added++
		}
	}
	return added
}

func SMembers(key string) []string {
	db.RLock()
	defer db.RUnlock()

	if set, exists := db.set[key]; exists {
		members := make([]string, 0, len(set))
		for member := range set {
			members = append(members, member)
		}
		return members
	}
	return []string{}
}

func SIsMember(key, member string) bool {
	db.RLock()
	defer db.RUnlock()

	if set, exists := db.set[key]; exists {
		_, exists := set[member]
		return exists
	}
	return false
}

func SRem(key string, members ...string) int {
	db.Lock()
	defer db.Unlock()

	removed := 0
	if set, exists := db.set[key]; exists {
		for _, member := range members {
			if _, exists := set[member]; exists {
				delete(set, member)
				removed++
			}
		}
		if len(set) == 0 {
			delete(db.set, key)
		}
	}
	return removed
}

// sorted sets

type SortedSet struct {
	members map[string]float64
	order   []string
}

var sortedSets = make(map[string]*SortedSet)

func ZAdd(key string, score float64, member string) string {
	if _, exists := sortedSets[key]; !exists {
		sortedSets[key] = &SortedSet{
			members: make(map[string]float64),
			order:   []string{},
		}
	}

	ss := sortedSets[key]
	ss.members[member] = score

	index := 0
	for i, m := range ss.order {
		if ss.members[m] > score {
			index = i
			break
		}
	}
	ss.order = append(ss.order[:index], append([]string{member}, ss.order[index:]...)...)
	return "OK"
}

func ZRange(key string, start, stop int) []string {
	if ss, exists := sortedSets[key]; exists {
		if start < 0 {
			start = 0
		}
		if stop >= len(ss.order) {
			stop = len(ss.order) - 1
		}
		return ss.order[start : stop+1]
	}
	return []string{}
}

func ZRem(key string, members ...string) int {
	if ss, exists := sortedSets[key]; exists {
		removedCount := 0
		for _, member := range members {
			if _, exists := ss.members[member]; exists {
				delete(ss.members, member)
				for i, m := range ss.order {
					if m == member {
						ss.order = append(ss.order[:i], ss.order[i+1:]...)
						break
					}
				}
				removedCount++
			}
		}
		return removedCount
	}
	return 0
}
