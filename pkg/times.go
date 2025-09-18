package pkg

import "time"

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (RealClock) Now() time.Time {
	return time.Now()
}

func GetAdjustedReleaseDay(releaseDate time.Time, now time.Time) int {
	releaseDay := releaseDate.YearDay()
	currentDay := now.YearDay()
	if isLeap(releaseDate) && !isLeap(now) && releaseDay >= 60 {
		return releaseDay - 1
	}
	if isLeap(now) && !isLeap(releaseDate) && currentDay >= 60 {
		return releaseDay + 1
	}
	return releaseDay
}

func isLeap(date time.Time) bool {
	year := date.Year()
	if year%400 == 0 {
		return true
	} else if year%100 == 0 {
		return false
	} else if year%4 == 0 {
		return true
	}
	return false
}
