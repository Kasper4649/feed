package time

import (
	"time"
)

func ParseTime(layout, value string) time.Time {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t, _ := time.ParseInLocation(layout, value, loc)
	return t
}
