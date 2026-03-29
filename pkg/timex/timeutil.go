package timex

import "time"

// ToUnixSeconds converts a time value to a Unix timestamp in seconds.
func ToUnixSeconds(value time.Time) int64 {
	return value.Unix()
}
