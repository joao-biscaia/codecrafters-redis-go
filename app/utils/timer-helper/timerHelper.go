package timerHelper

import (
	"time"

	constants "github.com/codecrafters-io/redis-starter-go/app/utils/consts"
)

type timeFunc map[string]time.Duration

func CreateTimeExpiry(durationMeasure string, duration int) time.Time {
	tFunc := timeFunc{
		constants.Milliseconds: time.Millisecond,
		constants.Seconds:      time.Second,
	}
	d := time.Duration(duration) * tFunc[durationMeasure]
	return time.Now().Add(d)

}
