package tool

import (
	"fmt"
	"strconv"
	"time"
)

// GenerateUniqId 生成唯一id
func GenerateUniqId(params string) string {
	return fmt.Sprint(Crc32(params + strconv.FormatInt(time.Now().Unix(), 10)))
}

// StringToInt64 string转int64
func StringToInt64(s string) int64 {
	d, _ := strconv.ParseInt(s, 10, 64)
	return d
}

// Int64ToString int64转string
func Int64ToString(i64 int64) string {
	return strconv.FormatInt(i64, 10)
}

// GetType 获取变量类型
func GetType(data interface{}) string {
	return fmt.Sprintf(`%T`, data)
}
