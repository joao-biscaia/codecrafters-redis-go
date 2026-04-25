package storage

import (
	"log"
	"sync"
	"time"

	timerHelper "github.com/codecrafters-io/redis-starter-go/app/utils/timer-helper"
)

var data sync.Map
var locks sync.Map

type entry struct {
	value     any
	expiresAt time.Time
}

func lockForKey(key string) *sync.Mutex {
	mu, _ := locks.LoadOrStore(key, &sync.Mutex{})
	return mu.(*sync.Mutex)
}

func Store(key string, value any) {
	log.Printf("storing key=%q value=%q", key, value)
	data.Store(key, entry{value: value})
}

func StoreWithExpiry(key string, value string, durationMeasure string, duration int) {
	e := entry{
		value:     value,
		expiresAt: timerHelper.CreateTimeExpiry(durationMeasure, duration),
	}
	data.Store(key, e)
}

func Get(key string) (string, bool) {
	v, ok := data.Load(key)
	if !ok {
		return "", false
	}
	e := v.(entry)
	if !e.expiresAt.IsZero() && time.Now().After(e.expiresAt) {
		data.Delete(key)
		return "", false
	}
	return e.value.(string), true
}

func Push(key string, args ...string) int {
	mu := lockForKey(key)
	mu.Lock()
	defer mu.Unlock()

	v, ok := data.LoadAndDelete(key)
	if ok {
		l := v.(entry).value.([]string)
		al := append(l, args...)
		data.Store(key, entry{value: al})
		return len(al)
	}
	nl := make([]string, len(args))
	copy(nl, args)
	data.Store(key, entry{value: nl})
	return len(nl)
}
