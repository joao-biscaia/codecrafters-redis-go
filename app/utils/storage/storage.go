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

func Get[T any](key string) (T, bool) {
	var zero T
	v, ok := data.Load(key)
	if !ok {
		return zero, false
	}
	e := v.(entry)
	if !e.expiresAt.IsZero() && time.Now().After(e.expiresAt) {
		data.Delete(key)
		return zero, false
	}
	val, ok := e.value.(T)
	return val, ok
}

func Push[T any](key string, r bool, args ...T) int {
	mu := lockForKey(key)
	mu.Lock()
	defer mu.Unlock()

	v, ok := data.LoadAndDelete(key)
	if ok {
		// in case values from args are different from l
		l, ok2 := v.(entry).value.([]T)
		if !ok2 {
			return -1
		}
		var al []T
		if r {
			al = append(l, args...)
		} else {
			al = append(args, l...)
		}
		data.Store(key, entry{value: al})
		return len(al)
	}
	nl := append([]T{}, args...)
	data.Store(key, entry{value: nl})
	return len(nl)
}
