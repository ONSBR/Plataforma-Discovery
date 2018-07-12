package util

import "time"

func Timestamp(t time.Time) int64 {
	return t.UTC().UnixNano() / int64(time.Millisecond)
}
