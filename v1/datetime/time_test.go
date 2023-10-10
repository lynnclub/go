package datetime

import (
	"fmt"
	"testing"
)

const timezone = "Asia/Shanghai"

// TestParseDateTime 解析日期时间
func TestParseDateTime(t *testing.T) {
	// ParseDateTime 解析日期时间
	goTime := ParseDateTime("2022", timezone)
	goTime2 := ParseDateTime("2022-11-08 15:30:02", timezone)
	if goTime.Unix() != 1640966400 {
		panic("datetime ParseDateTime error")
	}
	if goTime2.Unix() != 1667892602 {
		panic("datetime ParseDateTime error")
	}

	// ToUnix 转成时间戳
	timestamp := ToUnix("2022-11-08 15:30:02", timezone)
	if timestamp != 1667892602 {
		panic("datetime ToUnix error")
	}
}

// TestParseTimestamp 解析时间戳
func TestParseTimestamp(t *testing.T) {
	// ParseTimestamp 解析时间戳 秒级
	goTime := ParseTimestamp(1667895429, timezone)
	if goTime.Format(LayoutDateTime) != "2022-11-08 16:17:09" {
		panic("datetime ParseTimestamp error")
	}

	// ToAny 转成任意格式
	datetimeStr := ToAny(1667895429, timezone, LayoutDateTime)
	if datetimeStr != "2022-11-08 16:17:09" {
		panic("datetime ToAny error")
	}

	// ToDate 转成日期
	dateStr := ToDate(1667895429, timezone)
	if dateStr != "2022-11-08" {
		panic("datetime ToDate error")
	}

	// ToDateTime 转成日期时间
	datetimeStr = ToDateTime(1667895429, timezone)
	if datetimeStr != "2022-11-08 16:17:09" {
		panic("datetime ToDateTime error")
	}
}

// TestWeekAndNow
func TestWeekAndNow(t *testing.T) {
	// ToISOWeek 转成年周
	week := ToISOWeek(1667895429, timezone)
	if week != "2022_45" {
		panic("datetime ToISOWeek error")
	}

	// ToISOWeekByDate 转成年周
	week = ToISOWeekByDate("2022-11-08", timezone)
	if week != "2022_45" {
		panic("datetime ToISOWeekByDate error")
	}

	// Date 当前日期
	dateStr := Date(timezone)

	// DateTime 当前日期时间
	datetimeStr := DateTime(timezone)

	// ISOWeek 当前年周
	week = ISOWeek(timezone)

	fmt.Println(dateStr, datetimeStr, week)
}

// TestOther
func TestOther(t *testing.T) {
	// CheckTime 检查时间
	result := CheckTime(1667895429, "2022-11-01", "2023-01-01 00:00:00", timezone)
	if result != 1 {
		panic("datetime CheckTime error")
	}
}

// TestSingle
func TestSingle(t *testing.T) {
	Single.SetTimeZone(timezone)

	// ToUnix 转成时间戳
	timestamp := Single.ToUnix("2022-11-08 15:30:02")
	if timestamp != 1667892602 {
		panic("datetime Single.ToUnix error")
	}

	// ToAny 转成任意格式
	datetimeStr := Single.ToAny(1667895429, LayoutDateTime)
	if datetimeStr != "2022-11-08 16:17:09" {
		panic("datetime Single.ToAny error")
	}

	// ToDate 转成日期
	dateStr := Single.ToDate(1667895429)
	if dateStr != "2022-11-08" {
		panic("datetime Single.ToDate error")
	}

	// ToDateTime 转成日期时间
	datetimeStr = Single.ToDateTime(1667895429)
	if datetimeStr != "2022-11-08 16:17:09" {
		panic("datetime Single.ToDateTime error")
	}
}
