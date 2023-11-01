package time

import (
	"time"
)

func GetCurrentTimeByFormat(format string) string {
	l, _ := time.LoadLocation("Asia/Seoul")
	now := time.Now()
	currentTime := now.In(l)
	return currentTime.Format(format)
}
