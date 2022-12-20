package utils

import (
	"time"
)

// 日期转时间戳
func StrToTime(str string) int64 {
	stamp, _ := time.Parse("2006-01-02", str)
	return stamp.Unix()
}

// 时间戳格式化为日期字符串
func DateToStr(timestamp int64, tpl string) string {
	if tpl == "ymd" {
		tpl = "2006-01-02"
	} else if tpl == "ymdhi" {
		tpl = "2006-01-02 15:04"
	} else if tpl == "ymdhis" {
		tpl = "2006-01-02 15:04:05"
	} else if tpl == "his" {
		tpl = "15:04:05"
	} else if tpl == "hi" {
		tpl = "15:04"
	}
	return time.Unix(timestamp, 0).Format(tpl)
}
