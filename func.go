package tingle

import "time"

// FormatTimeToStr 格式化时间为字符串
func FormatTimeToStr(t *time.Time) string {
	return (*t).Format("2006-01-02 15:04:05")
}
