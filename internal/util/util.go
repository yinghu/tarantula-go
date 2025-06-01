package util

import (
	"crypto/rand"
	"encoding/base64"
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

func Key(size int16) []byte {
	key := make([]byte, size)
	rand.Read(key)
	return key
}

func KeyToBase64(key []byte) string {
	return base64.StdEncoding.EncodeToString(key)
}

func KeyFromBase64(key string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(key)
}
