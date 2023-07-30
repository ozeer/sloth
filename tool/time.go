package tool

import (
	"time"
)

// 日期模板
const timeLayout = "2006-01-02 15:04:05"

// TimestampToDate 时间戳转化为日期
func TimestampToDate(timestamp int64) string {
	return time.Unix(timestamp, 0).Format(timeLayout)
}

func CurrentDate() string {
	return time.Now().Format(timeLayout)
}
