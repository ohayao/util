package date

import (
	"regexp"
	"strings"
	"time"
)

const (
	YMD     = "yyyy/mm/dd"
	HMS     = "hh:MM:ss"
	HMSF    = "hh:MM:ss.fff"
	YMDHMS  = "yyy/mm/dd hh:MM:ss"
	YMDHMSF = "yyy/mm/dd hh:MM:ss.fff"
)

func replace(target, pattern, value string) string {
	reg := regexp.MustCompile(pattern)
	return reg.ReplaceAllString(target, value)
}

func format(t time.Time, layout string) string {
	d := strings.Split(t.Format("2006 01 02 15 04 05.000000"), " ")
	_ds := strings.Split(d[5], ".")
	d = append(d, _ds[1])
	d[5] = _ds[0]
	result := layout
	if strings.Contains(result, "y") {
		result = replace(result, "y{3,}", d[0])
		result = replace(result, "y+", d[0][:2])
	}
	if strings.Contains(result, "m") {
		result = replace(result, "m{2,}", d[1])
		result = replace(result, "m+", strings.TrimPrefix(d[1], "0"))
	}
	if strings.Contains(result, "d") {
		result = replace(result, "d{2,}", d[2])
		result = replace(result, "d+", strings.TrimPrefix(d[2], "0"))
	}
	if strings.Contains(result, "h") || strings.Contains(result, "H") {
		result = replace(result, "[h|H]{2,}", d[3])
		result = replace(result, "[h|H]+", strings.TrimPrefix(d[3], "0"))
	}
	if strings.Contains(result, "M") {
		result = replace(result, "M{2,}", d[4])
		result = replace(result, "M+", strings.TrimPrefix(d[4], "0"))
	}
	if strings.Contains(result, "s") || strings.Contains(result, "S") {
		result = replace(result, "[s|S]{2,}", d[5])
		result = replace(result, "[s|S]+", strings.TrimPrefix(d[5], "0"))
	}
	if strings.Contains(result, "f") || strings.Contains(result, "F") {
		result = replace(result, "[f|F]{4,}", d[6])
		result = replace(result, "[f|F]+", d[6][:3])
	}
	return result
}

// Now 格式化当前系统时间
func Now(fmt string) string {
	return format(time.Now(), fmt)
}

// 格式化时间
func Get(tm time.Time, fmt string) string {
	return format(tm, fmt)
}

// 格式化字符串时间
// tm 目标时间字符串；dft 目标时间字符串的格式
// 如果格式化失败，则返回空
func GetString(tm, dft, fmt string) string {
	if atime, err := time.Parse(dft, tm); err != nil {
		return ""
	} else {
		return format(atime, fmt)
	}
}

// 精确到微秒的时间戳格式化
func GetMicro(timestamp int64, fmt string) string {
	return format(time.UnixMicro(timestamp), fmt)
}

// 精确到毫米的时间戳格式化
func GetMilli(timestamp int64, fmt string) string {
	return format(time.UnixMilli(timestamp), fmt)
}
