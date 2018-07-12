package util

import "time"

func Timestamp(t time.Time) int64 {
	return t.UTC().UnixNano() / int64(time.Millisecond)
}

func ToISOString(t time.Time) string {
	return t.Format(time.RFC3339)
}

func TimeFromMilliTimestamp(timestamp int64) time.Time {
	return time.Unix(0, timestamp*int64(time.Millisecond))
}
