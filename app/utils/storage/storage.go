package storage

import "sync"

var data sync.Map

func Store(key string, value string) {
	data.Store(key, value)
}

func Get(key string) (string, bool) {
	v, ok := data.Load(key)
	if !ok {
		return "", false
	}
	return v.(string), true
}
