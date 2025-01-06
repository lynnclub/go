package datetime

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const timezone = "Asia/Shanghai"

func TestParse(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		input     any
		wantTime  time.Time
		wantError error
	}{
		{
			name:      "Time input",
			input:     now,
			wantTime:  now,
			wantError: nil,
		},
		{
			name:      "String input (valid date)",
			input:     "2025-01-06T15:04:05Z",
			wantTime:  time.Date(2025, 1, 6, 15, 4, 5, 0, time.UTC),
			wantError: nil,
		},
		{
			name:      "Invalid string input",
			input:     "invalid-date",
			wantTime:  time.Time{},
			wantError: errors.New("Could not find format for \"invalid-date\""),
		},
		{
			name:      "Integer input (Unix timestamp)",
			input:     int(1673029200),
			wantTime:  time.Unix(1673029200, 0),
			wantError: nil,
		},
		{
			name:      "Int64 input (Unix timestamp)",
			input:     int64(1673029200),
			wantTime:  time.Unix(1673029200, 0),
			wantError: nil,
		},
		{
			name:      "Float64 input (Unix timestamp with fractional seconds)",
			input:     float64(1673029200.123),
			wantTime:  time.Unix(1673029200, 122999906),
			wantError: nil,
		},
		{
			name:      "Primitive.DateTime input",
			input:     primitive.NewDateTimeFromTime(time.Date(2025, 1, 6, 15, 4, 5, 0, time.UTC)),
			wantTime:  time.Date(2025, 1, 6, 15, 4, 5, 0, time.UTC),
			wantError: nil,
		},
		{
			name:      "Primitive.Timestamp input",
			input:     primitive.Timestamp{T: uint32(1673029200), I: 0},
			wantTime:  time.Unix(1673029200, 0),
			wantError: nil,
		},
		{
			name:      "Unsupported map input",
			input:     primitive.M{"key": "value"},
			wantTime:  time.Time{},
			wantError: errors.New("not support primitive.M"),
		},
		{
			name:      "Unsupported slice type",
			input:     []int{1, 2, 3},
			wantTime:  time.Time{},
			wantError: errors.New("not support []int"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTime, gotErr := Parse(tt.input)
			if !gotTime.Equal(tt.wantTime) {
				t.Errorf("Parse(%v) got time = %v, want %v", tt.input, gotTime, tt.wantTime)
			}
			if !errorEquals(gotErr, tt.wantError) {
				t.Errorf("Parse(%v) got error = %v, want %v", tt.input, gotErr, tt.wantError)
			}
		})
	}
}

func errorEquals(err1, err2 error) bool {
	if err1 == nil || err2 == nil {
		return err1 == err2
	}
	return err1.Error() == err2.Error()
}

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
