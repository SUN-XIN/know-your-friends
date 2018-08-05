package main

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"

	"github.com/SUN-XIN/know-your-friends/helper"
	"github.com/SUN-XIN/know-your-friends/scylladb"
	"github.com/SUN-XIN/know-your-friends/types"
)

// calculate the total duration of user1 and user2 for 
// BestFriend (out of places) and MostSeen
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
