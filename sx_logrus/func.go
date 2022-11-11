package sx_logrus

import "time"

// 获取当前时间字符串
func now() string {
	return time.Now().Format(defaultTimestampFormat)
}
