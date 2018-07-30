package main

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"

	"github.com/SUN-XIN/know-your-friends/helper"
	"github.com/SUN-XIN/know-your-friends/scylladb"
	"github.com/SUN-XIN/know-your-friends/types"
)

func CalculDurationWithUser(dbSess *gocql.Session, ownerID, friendID string, inPlace bool) (totalDurationOut, totalDuration int32, err error) {
	// get stored SessionIntegrate from ScyllaDB
	days := helper.GetLastDays(time.Now())

	// TODO: GetMulti ?
	for _, d := range days {
		si := types.SessionIntegrate{
			UserIDOwner:   ownerID,
			UserIDFriend:  friendID,
			Day:           d,
			IsInSignPlace: inPlace,
		}

		err = scylladb.GetSessionIntegrate(dbSess, &si)
		if err != nil {
			if err == scylladb.ErrNotFound {
				err = nil
				continue
			}
			err = fmt.Errorf("Failed GetSessionIntegrate for %d: %+v", d, err)
			return
		}

		totalDuration = totalDuration + si.TotalDuration
		if !si.IsInSignPlace {
			totalDurationOut = totalDurationOut + si.TotalDuration
		}
	}

	return
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
