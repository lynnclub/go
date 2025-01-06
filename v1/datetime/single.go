package datetime

import "time"

var Single = &single{}

type single struct {
	timezone string
}

func (s *single) ParseAny(value any) (goTime time.Time, err error) {
	goTime, err = ParseAny(value)
	s.timezone = goTime.Location().String()

	return goTime, err
}

func (s *single) SetTimeZone(timezone string) {
	s.timezone = timezone
}

func (s *single) ParseDateTime(datetime string) time.Time {
	return ParseDateTime(datetime, s.timezone)
}

func (s *single) ToUnix(datetime string) int64 {
	return ToUnix(datetime, s.timezone)
}

func (s *single) ParseTimestamp(timestamp int64) time.Time {
	return ParseTimestamp(timestamp, s.timezone)
}

func (s *single) ToAny(timestamp int64, layout string) string {
	return ToAny(timestamp, s.timezone, layout)
}

func (s *single) ToDate(timestamp int64) string {
	return ToDate(timestamp, s.timezone)
}

func (s *single) ToDateTime(timestamp int64) string {
	return ToDateTime(timestamp, s.timezone)
}

func (s *single) ToISOWeek(timestamp int64) string {
	return ToISOWeek(timestamp, s.timezone)
}

func (s *single) ToISOWeekByDate(datetime string) string {
	return ToISOWeekByDate(datetime, s.timezone)
}

func (s *single) Unix() int64 {
	return Unix(s.timezone)
}

func (s *single) Date() string {
	return Date(s.timezone)
}

func (s *single) DateTime() string {
	return DateTime(s.timezone)
}

func (s *single) ISOWeek() string {
	return ISOWeek(s.timezone)
}

func (s *single) Any(layout string) string {
	return Any(s.timezone, layout)
}

func (s *single) CheckTime(timestamp int64, start string, end string) int {
	return CheckTime(timestamp, start, end, s.timezone)
}

func (s *single) CheckTimeNow(start string, end string) int {
	return CheckTimeNow(start, end, s.timezone)
}

func (s *single) ClickhouseDatatimeRange() (time.Time, time.Time) {
	return ClickhouseDatatimeRange()
}
