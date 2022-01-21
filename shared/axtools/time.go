package axtools

import "time"

type Interval uint32

const (
	Second Interval = iota
	Minute
	Hour
	Day
	Week
	Month
	Year
)

func GetStringTimeKey(t time.Time, interval Interval) string {
	return GetTimeKey(t, interval).Format("2006-01-02 15:04:05")
}

func GetStringTimeKeyNow(interval Interval) string {
	return GetStringTimeKey(time.Now(), interval)
}

func GetTimeKeyNow(interval Interval) time.Time {
	return GetTimeKey(time.Now(), interval)
}

func GetTimeKey(t time.Time, interval Interval) time.Time {
	switch interval {
	default:
		return t
	case Second:
		return t.Truncate(time.Second)
	case Minute:
		return t.Truncate(time.Minute)
	case Hour:
		return t.Truncate(time.Hour)
	case Day:
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	case Week:
		wd := t.Weekday()
		t = t.AddDate(0, 0, int(-wd))
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	case Month:
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	case Year:
		return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
	}
}
