package utilities

import (
	"strconv"
	"time"
)

func GetTimestamp() string {
	timestamp := (time.Now().UTC().UnixNano()) / 1000000
	return strconv.FormatInt(timestamp, 10)
}
