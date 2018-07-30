package main

import (
	"fmt"
	"log"
	"time"

	"github.com/SUN-XIN/know-your-friends/helper"
	"github.com/SUN-XIN/know-your-friends/scylladb"
	"github.com/SUN-XIN/know-your-friends/types"
)

func (s *server) CheckBestFriendAndMostSeen(ownerID, friendID string,
	lat, lng float64,
	startDate, endDate int64,
	resp *types.SessionReply) error {
	sessDay := helper.GetBeginningOfDay(time.Unix(endDate, 0))

	// if not found in cache -> get significant places from DB
	places := s.CheckPlaceInCache(ownerID)
	log.Printf("%d places in cache", len(places))

	dur := int32(endDate - startDate)
	isIn := false
	for _, p := range places {
		if p.IsIn(lat, lng) {
			isIn = true
			break
		}
	}
	log.Printf("this session is in some place (%v)", isIn)

	var err error
	/*
		// put session in ScyllaDB (only for backup)
		sd := types.SessionDetail{
			UserIDOwner:  ownerID,
			UserIDFriend: friendID,

			StartDate: startDate,
			EndDate:   endDate,

			Lat: lat,
			Lng: lng,

			IsInSignPlace: isIn,
		}
		err = scylladb.CreateSessionDetail(&sd)
		if err != nil {
			return fmt.Errorf("Failed PutSessionDetail: %+v", err)
		}
	*/

	// update SessionIntegrate in ScyllaDB
	si := types.SessionIntegrate{
		UserIDOwner:   ownerID,
		UserIDFriend:  friendID,
		Day:           sessDay,
		IsInSignPlace: isIn,
	}
	err = scylladb.UpdateSessionIntegrate(s.dbSession, &si, dur)
	if err != nil {
		return fmt.Errorf("Failed UpdateSessionIntegrate: %+v", err)
	}

	// need to re-calculate ?
	var topUser *types.TopUser
	var ok bool
	topUserInterf, existed := s.cacheUserTop.Get(ownerID)
	if !existed { // not in cache -> check in db
		log.Printf("not found TopUser in cache")
		topUser = &types.TopUser{
			OwnerID: ownerID,
			Day:     sessDay,
		}

		err = scylladb.GetTopUser(s.dbSession, topUser)
		switch err {
		case scylladb.ErrNotFound: // create top user
			log.Printf("not found TopUser in db -> create")
			topUser = types.NewTopUserByDuration(ownerID, friendID, startDate, endDate, isIn)

			// put in cache
			s.cacheUserTop.Add(ownerID, topUser)
			// put in db
			err = scylladb.PutTopUser(s.dbSession, topUser)
			if err != nil {
				return fmt.Errorf("Failed PutTopUser: %+v", err)
			}

			if !isIn {
				// Best Friend -> ok
				resp.BestFriend = friendID

				log.Printf("BestFriend: not in place -> write in response")
			}
			// Most seen -> ok
			resp.MostSeen = friendID

			log.Printf("MostSeen: write in response")

			return nil
		case nil: // found in db
		default:
			return fmt.Errorf("Failed GetTopUser: %+v", err)
		}
	}

	log.Printf("found TopUser in db -> check")
	topUser, ok = topUserInterf.(*types.TopUser)
	if !ok { // SHOULD NEVER HAPPEN
		return fmt.Errorf("SHOULD NEVER HAPPEN, Failed Conv to TopUser: %+v", topUserInterf)
	}

	// topUser existed
	var totalDurationOut, totalDuration int32

	// check "Most seen"
	if topUser.TopUserID == friendID && topUser.Day == sessDay {
		// update top user
		topUser.TopUserDuration = topUser.TopUserDuration + dur
		// put in db
		err = scylladb.PutTopUser(s.dbSession, topUser)
		// Most seen -> ok
		resp.MostSeen = friendID
		log.Printf("MostSeen: same top user -> write in response")
	} else {
		// calculate duration with friendID
		totalDurationOut, totalDuration, err = CalculDurationWithUser(s.dbSession, ownerID, friendID, isIn)
		if err != nil {
			return fmt.Errorf("Failed CalculDurationWithUser: %+v", err)
		}

		// update top user ?
		if totalDuration > topUser.TopUserDuration {
			topUser.TopUserID = friendID
			topUser.TopUserDuration = totalDuration

			// put in db
			err = scylladb.PutTopUser(s.dbSession, topUser)

			// Most seen -> ok
			resp.MostSeen = friendID

			log.Printf("MostSeen: diff top user -> update then write in response")
		} else {
			// keep stored top user
			// Most seen -> ok
			resp.MostSeen = topUser.TopUserID

			log.Printf("MostSeen: diff top user -> but keep top user then write in response")
		}
	}

	// check "Best Friend"
	if isIn { // keep the existed "Best Friend"
		// Best Friend -> ok
		resp.BestFriend = topUser.TopUserIDOutPlace

		log.Printf("BestFriend: in places -> write in response")

		return nil
	}

	// not in significant places -> check "Best Friend"
	if topUser.TopUserIDOutPlace == friendID && topUser.Day == sessDay {
		// update top user
		topUser.TopUserDurationOutPlace = topUser.TopUserDurationOutPlace + dur
		// put in db
		err = scylladb.PutTopUser(s.dbSession, topUser)
		// Best Friend -> ok
		resp.BestFriend = friendID

		log.Printf("BestFriend: in places and same top user -> write in response")

		return nil
	}

	if totalDurationOut == 0 {
		// calculate duration with friendID in significant places
		totalDurationOut, _, err = CalculDurationWithUser(s.dbSession, ownerID, friendID, isIn)
		if err != nil {
			return fmt.Errorf("Failed CalculDurationWithUser: %+v", err)
		}
	}

	// update top user ?
	if totalDurationOut > topUser.TopUserDurationOutPlace {
		topUser.TopUserIDOutPlace = friendID
		topUser.TopUserDurationOutPlace = totalDurationOut
		// put in db
		err = scylladb.PutTopUser(s.dbSession, topUser)
		// Best Friend -> ok
		resp.BestFriend = friendID

		log.Printf("BestFriend: in places and diff top user -> update then write in response")

		return nil
	}

	// keep top user
	// Best Friend -> ok
	resp.BestFriend = topUser.TopUserIDOutPlace

	log.Printf("BestFriend: in places and diff top user -> but keep top user then write in response")

	return nil
}

// check if user's Significant Places are in the cache
// and insert into cache if not
func (s *server) CheckPlaceInCache(userID string) []*types.SignificantPlace {
	places, stored := s.cacheUserPlaces.Get(userID)
	if !stored {
		placesOfUser := GetPlacesByID(userID)
		s.cacheUserPlaces.Add(userID, placesOfUser)

		return placesOfUser
	}

	// exist in cache
	placesOfUser, ok := places.([]*types.SignificantPlace)
	if !ok { // SHOULD NEVER HAPPEN
		log.Printf("SHOULD NEVER HAPPEN, Failed conv from cache value to place (%+v)", places)
		return nil
	}

	return placesOfUser
}
