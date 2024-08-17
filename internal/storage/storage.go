
package storage

var store = make(map[string]string)

func Set(key, value string) {
    store[key] = value
}

func Get(key string) string {
    return store[key]
}
