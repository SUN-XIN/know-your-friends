package helper

import (
	"time"
)

const (
	DEFAULT_ROLLING_DAYS = 7

	CRUSH_MIN_NIGHTS    = 3
	CRUSH_MIN_DURATION  = 21600 //6 * time.Hour
	CRUSH_DURATION_FROM = 22
	CRUSH_DURATION_TO   = 8
)

// return the timestamps for the beginning of each days (for last N days)
func GetLastDays(startDay time.Time) []int64 {
	res := make([]int64, 0, DEFAULT_ROLLING_DAYS)
	for i := 0; i < DEFAULT_ROLLING_DAYS; i++ {
		toDay := startDay.Add(-24 * time.Duration(int(i)) * time.Hour)

		res = append(res, GetBeginningOfDay(toDay))
	}

	return res
}

// return the timestamps for the beginning of the given day
func GetBeginningOfDay(d time.Time) int64 {
	// TODO: use user's time zone?
	return time.Date(d.Year(),
		d.Month(),
		d.Day(),
		0, 0, 0, 0,
		time.UTC).Unix()
}

func IsInPeriod(start, end int64) bool {
	// TODO: use user's time zone ?
	startDate := time.Unix(start, 0).UTC()
	endDate := time.Unix(end, 0).UTC()

	if startDate.Sub(endDate) > 24*time.Hour {
		return true
	}

	inSameDay := false
	if startDate.Day() == endDate.Day() {
		inSameDay = true
	}

	switch {
	case inSameDay &&
		startDate.Hour() < CRUSH_DURATION_TO &&
		endDate.Hour() < CRUSH_DURATION_TO:
		return true
	case !inSameDay &&
		startDate.Hour() > CRUSH_DURATION_FROM &&
		endDate.Hour() < CRUSH_DURATION_TO:
		return true
	default:
		return false
	}
}

/*
// fetch all session detail from scylladb
// then calculate sum of duration
func CalculDurationTotalOfDay(day int64) (durationInPlace, durationAll int32, err error) {
	var sessions []*types.SessionDetail
	sessions, err = scylladb.FetchAllSessionDetailOfDay(day)
	if err != nil {
		return
	}

	for _, s := range sessions {
		durationAll = durationAll + int32(s.EndDate-s.StartDate)
		if s.IsInSignPlace {
			durationInPlace = durationInPlace + int32(s.EndDate-s.StartDate)
		}
	}
	return
}
*/
