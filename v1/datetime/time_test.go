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
			gotTime, gotErr := ParseAny(tt.input)
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

// TestDiffInDays 测试天数差异计算
func TestDiffInDays(t *testing.T) {
	testCases := []struct {
		name     string
		t1       int64
		t2       int64
		expected int
	}{
		{"相同时间", 1667895429, 1667895429, 0},
		{"同一天不同时间", 1667865600, 1667895429, 0},
		{"相差一天", 1667779200, 1667865600, 1},
		{"相差多天", 1667779200, 1668124800, 4},
		{"跨月", 1667232000, 1669824000, 30},
		{"t1>t2", 1668124800, 1667779200, 4},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := DiffInDays(tc.t1, tc.t2)
			if result != tc.expected {
				t.Errorf("DiffInDays(%d, %d) = %d, 期望 %d", tc.t1, tc.t2, result, tc.expected)
			}
		})
	}
}

// TestGetWeekDay 测试获取本周起止日期
func TestGetWeekDay(t *testing.T) {
	start, end := GetWeekDay(timezone)

	if start == "" || end == "" {
		t.Error("GetWeekDay返回了空字符串")
	}

	// 验证格式
	startTime, err := time.Parse(LayoutDate, start)
	if err != nil {
		t.Errorf("起始日期格式错误: %v", err)
	}

	endTime, err := time.Parse(LayoutDate, end)
	if err != nil {
		t.Errorf("结束日期格式错误: %v", err)
	}

	// 结束日期应该大于起始日期
	if !endTime.After(startTime) {
		t.Errorf("结束日期应该大于起始日期: start=%s, end=%s", start, end)
	}

	// 起始日期应该是周一
	if startTime.Weekday() != time.Monday {
		t.Errorf("起始日期应该是周一，实际是%s", startTime.Weekday())
	}

	t.Logf("本周起止日期: %s 到 %s", start, end)
}

// TestClickhouseDatatimeRange 测试Clickhouse时间范围
func TestClickhouseDatatimeRange(t *testing.T) {
	start, end := ClickhouseDatatimeRange()

	expectedStart := "1970-01-01 00:00:00"
	expectedEnd := "2106-02-07 06:28:15"

	if start.Format(LayoutDateTime) != expectedStart {
		t.Errorf("起始时间错误，期望%s，实际%s", expectedStart, start.Format(LayoutDateTime))
	}

	if end.Format(LayoutDateTime) != expectedEnd {
		t.Errorf("结束时间错误，期望%s，实际%s", expectedEnd, end.Format(LayoutDateTime))
	}
}

// TestCheckTimeNow 测试检查当前时间
func TestCheckTimeNow(t *testing.T) {
	now := time.Now()

	// 测试当前时间在范围内
	yesterday := now.AddDate(0, 0, -1).Format(LayoutDateTime)
	tomorrow := now.AddDate(0, 0, 1).Format(LayoutDateTime)
	result := CheckTimeNow(yesterday, tomorrow, timezone)
	if result != 1 {
		t.Error("当前时间应该在范围内")
	}

	// 测试当前时间未开始
	future1 := now.AddDate(0, 0, 1).Format(LayoutDateTime)
	future2 := now.AddDate(0, 0, 2).Format(LayoutDateTime)
	result = CheckTimeNow(future1, future2, timezone)
	if result != 0 {
		t.Error("当前时间应该未开始")
	}

	// 测试当前时间已结束
	past1 := now.AddDate(0, 0, -2).Format(LayoutDateTime)
	past2 := now.AddDate(0, 0, -1).Format(LayoutDateTime)
	result = CheckTimeNow(past1, past2, timezone)
	if result != 2 {
		t.Error("当前时间应该已结束")
	}
}

// TestUnix 测试获取当前时间戳
func TestUnix(t *testing.T) {
	timestamp := Unix(timezone)
	now := time.Now().Unix()

	// 允许1秒误差
	if timestamp < now-1 || timestamp > now+1 {
		t.Errorf("Unix时间戳不准确: 期望约%d，实际%d", now, timestamp)
	}
}

// TestAny 测试自定义格式
func TestAny(t *testing.T) {
	result := Any(timezone, LayoutTime)

	// 验证格式是否正确
	_, err := time.Parse(LayoutTime, result)
	if err != nil {
		t.Errorf("Any返回的时间格式错误: %s", result)
	}
}

// TestSingleExtended 测试Single的更多方法
func TestSingleExtended(t *testing.T) {
	Single.SetTimeZone(timezone)

	// 测试ParseDateTime
	goTime := Single.ParseDateTime("2022-11-08 15:30:02")
	if goTime.Format(LayoutDateTime) != "2022-11-08 15:30:02" {
		t.Error("Single.ParseDateTime错误")
	}

	// 测试ParseTimestamp
	goTime = Single.ParseTimestamp(1667895429)
	if goTime.Format(LayoutDateTime) != "2022-11-08 16:17:09" {
		t.Error("Single.ParseTimestamp错误")
	}

	// 测试ToISOWeek
	week := Single.ToISOWeek(1667895429)
	if week != "2022_45" {
		t.Error("Single.ToISOWeek错误")
	}

	// 测试ToISOWeekByDate
	week = Single.ToISOWeekByDate("2022-11-08")
	if week != "2022_45" {
		t.Error("Single.ToISOWeekByDate错误")
	}

	// 测试Unix
	timestamp := Single.Unix()
	now := time.Now().Unix()
	if timestamp < now-1 || timestamp > now+1 {
		t.Error("Single.Unix错误")
	}

	// 测试Date
	dateStr := Single.Date()
	if dateStr == "" {
		t.Error("Single.Date错误")
	}

	// 测试DateTime
	datetimeStr := Single.DateTime()
	if datetimeStr == "" {
		t.Error("Single.DateTime错误")
	}

	// 测试ISOWeek
	week = Single.ISOWeek()
	if week == "" {
		t.Error("Single.ISOWeek错误")
	}

	// 测试Any
	result := Single.Any(LayoutTime)
	_, err := time.Parse(LayoutTime, result)
	if err != nil {
		t.Error("Single.Any错误")
	}

	// 测试CheckTime
	result2 := Single.CheckTime(1667895429, "2022-11-01", "2023-01-01")
	if result2 != 1 {
		t.Error("Single.CheckTime错误")
	}

	// 测试CheckTimeNow
	now2 := time.Now()
	yesterday := now2.AddDate(0, 0, -1).Format(LayoutDateTime)
	tomorrow := now2.AddDate(0, 0, 1).Format(LayoutDateTime)
	result2 = Single.CheckTimeNow(yesterday, tomorrow)
	if result2 != 1 {
		t.Error("Single.CheckTimeNow错误")
	}

	// 测试ClickhouseDatatimeRange
	start, end := Single.ClickhouseDatatimeRange()
	if start.IsZero() || end.IsZero() {
		t.Error("Single.ClickhouseDatatimeRange错误")
	}
}

// TestParseDateTimeEdgeCases 测试ParseDateTime边界情况
func TestParseDateTimeEdgeCases(t *testing.T) {
	// 空字符串应该返回当前时间
	result := ParseDateTime("", timezone)
	now := time.Now()
	if result.Year() != now.Year() {
		t.Error("空字符串应该返回当前时间")
	}

	// 不同长度的日期字符串
	testCases := []struct {
		input    string
		expected string
	}{
		{"2022", "2022-01-01 00:00:00"},
		{"2022-11", "2022-11-01 00:00:00"},
		{"2022-11-08", "2022-11-08 00:00:00"},
		{"2022-11-08 15", "2022-11-08 15:00:00"},
		{"2022-11-08 15:30", "2022-11-08 15:30:00"},
		{"2022-11-08 15:30:02", "2022-11-08 15:30:02"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := ParseDateTime(tc.input, timezone)
			if result.Format(LayoutDateTime) != tc.expected {
				t.Errorf("ParseDateTime(%q) = %s, 期望 %s", tc.input, result.Format(LayoutDateTime), tc.expected)
			}
		})
	}
}

// TestParseTimestampEdgeCases 测试ParseTimestamp边界情况
func TestParseTimestampEdgeCases(t *testing.T) {
	// 0或负数时间戳应该返回当前时间
	result := ParseTimestamp(0, timezone)
	now := time.Now()
	if result.Year() != now.Year() {
		t.Error("0时间戳应该返回当前时间")
	}

	result = ParseTimestamp(-1, timezone)
	if result.Year() != now.Year() {
		t.Error("负数时间戳应该返回当前时间")
	}

	// 无效时区应该使用UTC
	result = ParseTimestamp(1667895429, "Invalid/Timezone")
	if result.Unix() != 1667895429 {
		t.Error("无效时区应该仍然能解析时间戳")
	}
}

// TestByteSliceInput 测试[]byte输入
func TestByteSliceInput(t *testing.T) {
	input := []byte("2025-01-06T15:04:05Z")
	goTime, err := ParseAny(input)

	if err != nil {
		t.Errorf("ParseAny([]byte) 返回错误: %v", err)
	}

	expected := time.Date(2025, 1, 6, 15, 4, 5, 0, time.UTC)
	if !goTime.Equal(expected) {
		t.Errorf("ParseAny([]byte) = %v, 期望 %v", goTime, expected)
	}
}

// TestSingleParseAny 测试Single的ParseAny方法
func TestSingleParseAny(t *testing.T) {
	Single.SetTimeZone(timezone)

	// 测试有效输入
	input := "2025-01-06T15:04:05Z"
	goTime, err := Single.ParseAny(input)

	if err != nil {
		t.Errorf("Single.ParseAny返回错误: %v", err)
	}

	if goTime.Location().String() != timezone {
		t.Errorf("时区应该是%s，实际是%s", timezone, goTime.Location().String())
	}

	// 测试无效时区
	Single.SetTimeZone("Invalid/Timezone")
	_, err = Single.ParseAny(input)
	if err == nil {
		t.Error("无效时区应该返回错误")
	}

	// 恢复时区
	Single.SetTimeZone(timezone)
}
