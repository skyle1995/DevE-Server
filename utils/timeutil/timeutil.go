package timeutil

import (
	"fmt"
	"time"
)

const (
	// 常用时间格式
	DateFormat     = "2006-01-02"
	TimeFormat     = "15:04:05"
	DateTimeFormat = "2006-01-02 15:04:05"
	CompactFormat  = "20060102150405"
)

// FormatTime 格式化时间
// t: 时间
// layout: 格式，可以使用预定义的常量，如DateFormat, TimeFormat, DateTimeFormat
func FormatTime(t time.Time, layout string) string {
	return t.Format(layout)
}

// ParseTime 解析时间字符串
// timeStr: 时间字符串
// layout: 格式，可以使用预定义的常量，如DateFormat, TimeFormat, DateTimeFormat
func ParseTime(timeStr, layout string) (time.Time, error) {
	return time.Parse(layout, timeStr)
}

// Now 获取当前时间
func Now() time.Time {
	return time.Now()
}

// NowStr 获取当前时间字符串
// layout: 格式，可以使用预定义的常量，如DateFormat, TimeFormat, DateTimeFormat
func NowStr(layout string) string {
	return FormatTime(Now(), layout)
}

// AddDuration 增加时间
// t: 时间
// d: 时间间隔
func AddDuration(t time.Time, d time.Duration) time.Time {
	return t.Add(d)
}

// AddSeconds 增加秒数
// t: 时间
// seconds: 秒数
func AddSeconds(t time.Time, seconds int) time.Time {
	return t.Add(time.Duration(seconds) * time.Second)
}

// AddMinutes 增加分钟数
// t: 时间
// minutes: 分钟数
func AddMinutes(t time.Time, minutes int) time.Time {
	return t.Add(time.Duration(minutes) * time.Minute)
}

// AddHours 增加小时数
// t: 时间
// hours: 小时数
func AddHours(t time.Time, hours int) time.Time {
	return t.Add(time.Duration(hours) * time.Hour)
}

// AddDays 增加天数
// t: 时间
// days: 天数
func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

// AddMonths 增加月数
// t: 时间
// months: 月数
func AddMonths(t time.Time, months int) time.Time {
	return t.AddDate(0, months, 0)
}

// AddYears 增加年数
// t: 时间
// years: 年数
func AddYears(t time.Time, years int) time.Time {
	return t.AddDate(years, 0, 0)
}

// DiffSeconds 计算两个时间的秒数差
// t1, t2: 两个时间
// 返回t1-t2的秒数差
func DiffSeconds(t1, t2 time.Time) int64 {
	return t1.Unix() - t2.Unix()
}

// DiffDays 计算两个时间的天数差
// t1, t2: 两个时间
// 返回t1-t2的天数差
func DiffDays(t1, t2 time.Time) int {
	// 将时间调整到当天的0点
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, t1.Location())
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, t2.Location())
	return int(t1.Sub(t2).Hours() / 24)
}

// IsLeapYear 判断是否为闰年
// year: 年份
func IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// GetMonthDays 获取某月的天数
// year: 年份
// month: 月份
func GetMonthDays(year, month int) int {
	if month < 1 || month > 12 {
		return 0
	}

	days := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	if month == 2 && IsLeapYear(year) {
		return 29
	}
	return days[month-1]
}

// FormatDuration 格式化时间间隔
// d: 时间间隔
// 返回格式化后的字符串，如：1天2小时3分钟4秒
func FormatDuration(d time.Duration) string {
	totalSeconds := int(d.Seconds())
	days := totalSeconds / (24 * 3600)
	hours := (totalSeconds % (24 * 3600)) / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	result := ""
	if days > 0 {
		result += fmt.Sprintf("%d天", days)
	}
	if hours > 0 {
		result += fmt.Sprintf("%d小时", hours)
	}
	if minutes > 0 {
		result += fmt.Sprintf("%d分钟", minutes)
	}
	if seconds > 0 || result == "" {
		result += fmt.Sprintf("%d秒", seconds)
	}

	return result
}

// DaysBetween 计算两个时间之间的天数
// now: 当前时间
// future: 未来时间
// 返回未来时间与当前时间之间的天数差
func DaysBetween(now, future time.Time) int {
	// 将时间调整到当天的0点
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	future = time.Date(future.Year(), future.Month(), future.Day(), 0, 0, 0, 0, future.Location())

	// 计算天数差
	duration := future.Sub(now)
	return int(duration.Hours() / 24)
}
