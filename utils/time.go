package utils

import "time"

func NowTruncated() time.Time {
	return time.Now().Truncate(time.Millisecond)
}
