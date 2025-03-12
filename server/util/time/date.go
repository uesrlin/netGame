package time

import (
	"time"
)

/**
判断当前时间是否在指定时间范围内
*/
func BetweenTime(start, end time.Time) bool {
	if time.Now().Unix() >= start.Unix() && time.Now().Unix() <= end.Unix() {
		return true
	}
	return false
}

/**
判断两个时间的差距
*/
func DistanceTime(t1, t2 time.Time) int64 {
	dim := t1.Unix() - t2.Unix()
	if dim >= 0 {
		return dim
	} else {
		return -dim
	}
}
