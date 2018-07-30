package main

import (
	"fmt"
	"log"
	"time"

	"github.com/SUN-XIN/know-your-friends/helper"
	"github.com/SUN-XIN/know-your-friends/scylladb"
	"github.com/SUN-XIN/know-your-friends/types"
)

func (s *server) CheckCrush(sess *types.SessionRequest,
	resp *types.SessionReply) error {
	// more than 6h ?
	if sess.EndDate-sess.StartDate < helper.CRUSH_MIN_DURATION {
		log.Printf("CheckCrush: less than 6h")
		return nil
	}

	// match period night
	if !helper.IsInPeriod(sess.StartDate, sess.EndDate) {
		log.Printf("CheckCrush: out of period")
		return nil
	}

	// if not found in cache -> get home from DB
	homeUser1 := s.CheckHomeInCache(sess.UserID1)
	homeUser2 := s.CheckHomeInCache(sess.UserID2)

	// in home ?
	if (homeUser1 != nil && homeUser1.IsIn(sess.Latitude, sess.Longitude)) ||
		(homeUser2 != nil && homeUser2.IsIn(sess.Latitude, sess.Longitude)) {
		// put in db
		sr, err := scylladb.GetAndUpdateSessionCrush(s.dbSession,
			sess.UserID1,
			sess.UserID2,
			helper.GetBeginningOfDay(time.Unix(sess.EndDate, 0)))
		if err != nil {
			return fmt.Errorf("Failed GetAndUpdateSessionCrush: %+v", err)
		}

		// TODO: use TopUser (cache)

		// check Crush
		nights, err := scylladb.CountNights(s.dbSession, sr)
		if err != nil {
			return fmt.Errorf("Failed CountNights: %+v", err)
		}

		if nights < helper.CRUSH_MIN_NIGHTS {
			log.Printf("CheckCrush: only %d nights", nights)
			return nil
		}

		resp.Crush = sr.FriendsIDs
		log.Printf("CheckCrush: write in response")

		return nil
	}

	log.Printf("CheckCrush: not in home")
	return nil
}

// check if user's Homes are in the cache
// and insert into cache if not
func (s *server) CheckHomeInCache(userID string) *types.SignificantPlace {
	placesInterf, stored := s.cacheUserPlaces.Get(userID)
	if !stored {
		// get from db
		placesOfUser := GetPlacesByID(userID)
		s.cacheUserPlaces.Add(userID, placesOfUser)

		// have home ?
		for _, p := range placesOfUser {
			if p.Name == types.PLACE_NAME_HOME {
				return p
			}
		}
		return nil
	}

	// found in cache
	places, ok := placesInterf.([]*types.SignificantPlace)
	if !ok { // SHOULD NEVER HAPPEN
		log.Printf("SHOULD NEVER HAPPEN, Failed conv to []*SignificantPlace")
		return nil
	}

	// exist in cache -> have home ?
	for _, p := range places {
		if p.Name == types.PLACE_NAME_HOME {
			return p
		}
	}

	return nil
}
