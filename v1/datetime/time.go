package datetime

import (
	"errors"
	"math"
	"reflect"
	"strconv"
	"time"

	"github.com/araddon/dateparse"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const LayoutTime = "15:04:05"
const LayoutDate = "2006-01-02"
const LayoutDateHour = "2006-01-02 15"
const LayoutDateTime = "2006-01-02 15:04:05"
const LayoutDateTimeZone = "2006-01-02 15:04:05 -0700 MST"
const LayoutDateTimeZoneT = "2006-01-02T15:04:05.999999-07:00"
const LayoutMonth = "2006-01"

// ParseAny 解析时间
func ParseAny(value any) (goTime time.Time, err error) {
	switch v := value.(type) {
	case time.Time:
		goTime = v
		err = nil
	case []byte:
		goTime, err = dateparse.ParseAny(string(v))
	case string:
		goTime, err = dateparse.ParseAny(v)
	case int:
		goTime = time.Unix(int64(v), 0)
		err = nil
	case int64:
		goTime = time.Unix(v, 0)
		err = nil
	case float64:
		sec, dec := math.Modf(v)
		goTime = time.Unix(int64(sec), int64(dec*(1e9)))
		err = nil
	case primitive.DateTime:
		goTime = v.Time()
		err = nil
	case primitive.Timestamp:
		goTime = time.Unix(int64(v.T), int64(v.I))
		err = nil
	default:
		err = errors.New("not support " + reflect.TypeOf(v).String())
	}

	return goTime, err
}

// ParseDateTime 解析日期时间
func ParseDateTime(datetime string, timezone string) time.Time {
	local, _ := time.LoadLocation(timezone)
	length := len(datetime)

	var goTime time.Time
	if local == nil {
		if length > 0 {
			goTime, _ = time.Parse(LayoutDateTime[0:length], datetime)
		} else {
			goTime = time.Now()
		}
	} else {
		if length > 0 {
			goTime, _ = time.ParseInLocation(LayoutDateTime[0:length], datetime, local)
		} else {
			goTime = time.Now().In(local)
		}
	}

	return goTime
}

// ToUnix 转成时间戳
func ToUnix(datetime string, timezone string) int64 {
	return ParseDateTime(datetime, timezone).Unix()
}

// ParseTimestamp 解析时间戳 秒级
func ParseTimestamp(timestamp int64, timezone string) time.Time {
	var goTime time.Time
	if timestamp > 0 {
		goTime = time.Unix(timestamp, 0)
	} else {
		goTime = time.Now()
	}

	local, _ := time.LoadLocation(timezone)
	if local == nil {
		return goTime
	} else {
		return goTime.In(local)
	}
}

// ToAny 转成任意格式
func ToAny(timestamp int64, timezone string, layout string) string {
	return ParseTimestamp(timestamp, timezone).Format(layout)
}

// ToDate 转成日期
func ToDate(timestamp int64, timezone string) string {
	return ToAny(timestamp, timezone, LayoutDate)
}

// ToDateTime 转成日期时间
func ToDateTime(timestamp int64, timezone string) string {
	return ToAny(timestamp, timezone, LayoutDateTime)
}

// ToISOWeek 转成年周 例如2020_5
func ToISOWeek(timestamp int64, timezone string) string {
	year, week := ParseTimestamp(timestamp, timezone).ISOWeek()

	return strconv.Itoa(year) + "_" + strconv.Itoa(week)
}

// ToISOWeekByDate 转成年周 例如2020_5
func ToISOWeekByDate(datetime string, timezone string) string {
	year, week := ParseDateTime(datetime, timezone).ISOWeek()

	return strconv.Itoa(year) + "_" + strconv.Itoa(week)
}

// Unix 当前时间戳
func Unix(timezone string) int64 {
	return ToUnix("", timezone)
}

// Date 当前日期
func Date(timezone string) string {
	return ToAny(0, timezone, LayoutDate)
}

// DateTime 当前日期时间
func DateTime(timezone string) string {
	return ToAny(0, timezone, LayoutDateTime)
}

// ISOWeek 当前年周 例如2020_5
func ISOWeek(timezone string) string {
	return ToISOWeek(0, timezone)
}

// Any 当前任意格式
func Any(timezone, layout string) string {
	return ToAny(0, timezone, layout)
}

// CheckTime 检查时间 0未开始、1正常、2已结束，半闭合区间[start, end)
func CheckTime(timestamp int64, start string, end string, timezone string) int {
	startTime := ToUnix(start, timezone)
	endTime := ToUnix(end, timezone)
	if timestamp < startTime {
		return 0
	}
	if timestamp >= endTime {
		return 2
	}

	return 1
}

// CheckTimeNow 检查当前时间 0未开始、1正常、2已结束，半闭合区间[start, end)
func CheckTimeNow(start string, end string, timezone string) int {
	return CheckTime(time.Now().Unix(), start, end, timezone)
}

// DiffInDays 时间戳相差天数
func DiffInDays(t1, t2 int64) int {
	if t1 == t2 {
		return 0
	}

	if t1 > t2 {
		t1, t2 = t2, t1
	}

	days := 0
	secDiff := t2 - t1
	if secDiff > 86400 {
		tmpDays := int(secDiff / 86400)
		t1 += int64(tmpDays) * 86400
		days += tmpDays
	}

	start := time.Unix(t1, 0)
	end := time.Unix(t2, 0)
	if start.Format(LayoutDate) != end.Format(LayoutDate) {
		days += 1
	}

	return days
}

// GetWeekDay 获取本周起止日期
func GetWeekDay(timezone string) (string, string) {
	loc, _ := time.LoadLocation(timezone)
	now := time.Now()

	offset := int(time.Monday - now.Weekday())
	// 周日特殊判断 因为time.Monday = 0
	if offset > 0 {
		offset = -6
	}

	lastOffset := int(time.Saturday - now.Weekday())
	// 周日特殊判断 因为time.Monday = 0
	if lastOffset == 6 {
		lastOffset = -1
	}

	firstOfWeek := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, offset)
	lastOfWeeK := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, lastOffset+1)
	start := firstOfWeek.Unix()
	end := lastOfWeeK.Unix()

	return time.Unix(start, 0).Format(LayoutDate), time.Unix(end, 0).Format(LayoutDate)
}

// clickhouse datetime类型底层以时间戳存储，不包含时区
func ClickhouseDatatimeRange() (time.Time, time.Time) {
	start, _ := time.Parse(LayoutDateTime, "1970-01-01 00:00:00")
	end, _ := time.Parse(LayoutDateTime, "2106-02-07 06:28:15")

	return start, end
}
