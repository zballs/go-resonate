package util

import "strconv"
import "time"

func TimeString() string {
	return time.Now().String()
}

func DateString() string {
	return ToTheDay(TimeString())
}

func ParseTimeString(timestr string) time.Time {
	yr, _ := strconv.Atoi(timestr[:4])
	mo, _ := strconv.Atoi(timestr[5:7])
	d, _ := strconv.Atoi(timestr[8:10])
	hr, _ := strconv.Atoi(timestr[11:13])
	min, _ := strconv.Atoi(timestr[14:16])
	sec, _ := strconv.Atoi(timestr[17:19])
	return time.Date(yr, time.Month(mo), d, hr, min, sec, 0, time.Local)
}

func ParseDateString(datestr string) time.Time {
	yr, _ := strconv.Atoi(datestr[:4])
	mo, _ := strconv.Atoi(datestr[5:7])
	d, _ := strconv.Atoi(datestr[8:10])
	return time.Date(yr, time.Month(mo), d, 0, 0, 0, 0, time.Local)
}

func ToTheDay(timestr string) string {
	return timestr[:10]
}

var months = map[string]int{
	"Jan": 1,
	"Feb": 2,
	"Mar": 3,
	"Apr": 4,
	"May": 5,
	"Jun": 6,
	"Jul": 7,
	"Aug": 8,
	"Sep": 9,
	"Oct": 10,
	"Nov": 11,
	"Dec": 12,
}
