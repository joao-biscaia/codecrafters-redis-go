package storage

import (
	"sync"
	"time"

	timerHelper "github.com/codecrafters-io/redis-starter-go/app/utils/timer-helper"
)

var data sync.Map

type entry struct {
	value     string
	expiresAt time.Time
}

func Store(key string, value string) {
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
	return e.value, true
}
