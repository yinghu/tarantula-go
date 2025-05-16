package util

import (
	"encoding/json"
	"hash/fnv"
	"time"
)

func EpochMillisecondsFromMidnight(year int, month int, day int) int64 {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return t.UnixMilli()
}

func Partition(key []byte, pnumber uint32) uint32 {
	hash := fnv.New32()
	hash.Write(key)
	return hash.Sum32() % pnumber
}

func ToJson(obj any) []byte {
	data, err := json.Marshal(obj)
	if err != nil {
		return []byte("{}")
	}
	return data
}
