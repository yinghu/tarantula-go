package util

import "time"


func EpochMillisecondsFromMidnight(year int, month int, day int) int64 {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return t.UnixMilli()
}
