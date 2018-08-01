package main

import (
	"fmt"
	"log"
	"time"

	"github.com/SUN-XIN/know-your-friends/helper"
	"github.com/SUN-XIN/know-your-friends/scylladb"
	"github.com/SUN-XIN/know-your-friends/types"
)

func (s *server) GetCrush(ownerID string, day int64, resp *types.UserFriendsReply) error {
	// already processed today ?
	tu := types.TopUser{
		OwnerID: ownerID,
		Day:     day,
	}
	err := scylladb.GetTopUser(s.dbSession, &tu)
	switch {
	case (err == nil && len(tu.CrushFriendIDs) > 0):
		resp.Crush = tu.CrushFriendIDs
		return nil
	case err == scylladb.ErrNotFound ||
		(err == nil && len(tu.CrushFriendIDs) == 0):
		resp.Crush, err = s.CheckAndPutCrush(&types.SessionCrush{
			UserIDOwner: ownerID,
			Day:         day,
		})
		if err != nil {
			return fmt.Errorf("Failed CheckAndPutCrush: %+v", err)
		}

		return nil
	default:
		return err
	}
}

func (s *server) CheckCrush(ownerID, friendID string,
	startDate, endDate int64,
	lat, lng float64) error {
	// more than 6h ?
	if endDate-startDate < helper.CRUSH_MIN_DURATION {
		log.Printf("CheckCrush: less than 6h")
		return nil
	}

	// match period night
	if !helper.IsInPeriod(startDate, endDate) {
		log.Printf("CheckCrush: out of period")
		return nil
	}

	// if not found in cache -> get home from DB
	homeUser1 := s.CheckHomeInCache(ownerID)
	homeUser2 := s.CheckHomeInCache(friendID)

	// in home ?
	if (homeUser1 != nil && homeUser1.IsIn(lat, lng)) ||
		(homeUser2 != nil && homeUser2.IsIn(lat, lng)) {
		// put in db
		sr, err := scylladb.GetAndUpdateSessionCrush(s.dbSession,
			ownerID,
			friendID,
			helper.GetBeginningOfDay(time.Unix(endDate, 0)))
		if err != nil {
			return fmt.Errorf("Failed GetAndUpdateSessionCrush: %+v", err)
		}

		// TODO: use TopUser (cache)

		// check Crush
		_, err = s.CheckAndPutCrush(sr)
		if err != nil {
			return fmt.Errorf("Failed CheckAndPutCrush: %+v", err)
		}

		return nil
	}

	log.Printf("CheckCrush: not in home")
	return nil
}

func (s *server) CheckAndPutCrush(sr *types.SessionCrush) ([]string, error) {
	fIDs, err := scylladb.CountNights(s.dbSession, sr)
	if err != nil {
		return nil, fmt.Errorf("Failed CountNights: %+v", err)
	}

	if len(fIDs) == 0 {
		return nil, nil
	}

	_, err = scylladb.UpdateTopUserCrush(s.dbSession,
		sr.UserIDOwner, fIDs, sr.Day)
	if err != nil {
		return nil, fmt.Errorf("Failed UpdateTopUserCrush: %+v", err)
	}
	log.Printf("CheckCrush: write in db")
	return fIDs, nil
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
