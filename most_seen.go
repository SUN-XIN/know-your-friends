package main

import (
	"fmt"
	"log"
	"time"

	"github.com/SUN-XIN/know-your-friends/helper"
	"github.com/SUN-XIN/know-your-friends/scylladb"
	"github.com/SUN-XIN/know-your-friends/types"
)

func (s *server) GetBestFriendAndMostSeen(ownerID string, day int64, resp *types.UserFriendsReply) error {
	// in cache ?
	var topUser *types.TopUser
	var ok bool
	topUserInterf, existed := s.cacheUserTop.Get(ownerID)
	if existed {
		topUser, ok = topUserInterf.(*types.TopUser)
		if !ok { // SHOULD NEVER HAPPEN
			return fmt.Errorf("SHOULD NEVER HAPPEN, Failed Conv to TopUser: %+v", topUserInterf)
		}
		log.Printf("Find in cache")
		resp.BestFriend = topUser.TopUserIDOutPlace
		resp.MostSeen = topUser.TopUserID
		return nil
	}

	log.Printf("Not in cache")
	days := helper.GetLastDays(time.Unix(day, 0).UTC())
	bestFriendsDur := make(map[string]int32, DEFAULT_FRIENDS_PER_DAY)
	mostSeenDur := make(map[string]int32, DEFAULT_FRIENDS_PER_DAY)
	var dur int32
	var got bool
	for _, d := range days {
		users, err := scylladb.QuerySessionIntegrate(s.dbSession, ownerID, d)
		if err != nil {
			return fmt.Errorf("Failed QuerySessionIntegrate: %+v", err)
		}

		log.Printf("%d friends for day %d", len(users), d)
		for _, u := range users {
			log.Printf("friend %v", *u)

			dur, got = mostSeenDur[u.UserIDFriend]
			if got {
				dur = dur + u.TotalDuration
			}
			mostSeenDur[u.UserIDFriend] = dur

			if !u.IsInSignPlace { // most seen
				dur, got = bestFriendsDur[u.UserIDFriend]
				if got {
					dur = dur + u.TotalDuration
				}
				bestFriendsDur[u.UserIDFriend] = dur
			}
		}
	}

	log.Printf("%d friends, %d friends out place,", len(mostSeenDur), len(bestFriendsDur))
	log.Printf("sort before: %v", mostSeenDur)

	resBestFriends := helper.SortMap(bestFriendsDur)
	resMostSeen := helper.SortMap(mostSeenDur)

	log.Printf("sort res: %v", resMostSeen)

	resp.BestFriend = resBestFriends[0].Key
	resp.MostSeen = resMostSeen[0].Key

	return nil
}

func (s *server) CheckBestFriendAndMostSeen(ownerID, friendID string,
	lat, lng float64,
	startDate, endDate int64) error {
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

	// update SessionIntegrate in ScyllaDB
	si := types.SessionIntegrate{
		UserIDOwner:   ownerID,
		UserIDFriend:  friendID,
		Day:           sessDay,
		IsInSignPlace: isIn,
	}
	err := scylladb.UpdateSessionIntegrate(s.dbSession, &si, dur)
	if err != nil {
		return fmt.Errorf("Failed UpdateSessionIntegrate: %+v", err)
	}

	cacheKey := fmt.Sprintf("%s-%d", ownerID, sessDay)

	// need to re-calculate ?
	var topUser *types.TopUser
	var ok bool
	topUserInterf, existed := s.cacheUserTop.Get(cacheKey)
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
			s.cacheUserTop.Add(cacheKey, topUser)
			// put in db
			err = scylladb.PutTopUser(s.dbSession, topUser)
			if err != nil {
				return fmt.Errorf("Failed PutTopUser: %+v", err)
			}

			if !isIn {
				log.Printf("%s is BestFriend of %s for day %d", friendID, ownerID, sessDay)
			}
			log.Printf("%s is MostSeen of %s for day %d", friendID, ownerID, sessDay)

			return nil
		case nil: // found in db
		default:
			return fmt.Errorf("Failed GetTopUser for %s %s %d: %+v", friendID, ownerID, sessDay, err)
		}
	}

	if existed {
		topUser, ok = topUserInterf.(*types.TopUser)
		if !ok { // SHOULD NEVER HAPPEN
			return fmt.Errorf("SHOULD NEVER HAPPEN, Failed Conv to TopUser: %+v", topUserInterf)
		}
	}
	log.Printf("found TopUser in db/cache for %s %s %d: %+v -> check", friendID, ownerID, sessDay)

	// topUser existed
	var totalDurationOut, totalDuration int32

	// check "Most seen"
	if topUser.TopUserID == friendID {
		// update top user
		topUser.TopUserDuration = topUser.TopUserDuration + dur
		// put in db
		err = scylladb.PutTopUser(s.dbSession, topUser)
		if err != nil {
			return fmt.Errorf("Failed PutTopUser for %s %s %d: %+v", friendID, ownerID, sessDay, err)
		}
		// Most seen -> ok
		log.Printf("%s is MostSeen (same top user) of %s for day %d", friendID, ownerID, sessDay)
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
			if err != nil {
				return fmt.Errorf("Failed PutTopUser for %s %s %d: %+v", friendID, ownerID, sessDay, err)
			}

			// Most seen -> ok
			log.Printf("%s is MostSeen (diff top user, update) of %s for day %d", friendID, ownerID, sessDay)
		} else {
			// keep stored top user
			log.Printf("%s is MostSeen (diff top user, keep) of %s for day %d", friendID, ownerID, sessDay)
		}
	}

	// check "Best Friend"
	if isIn { // keep the existed "Best Friend"
		log.Printf("%s is BestFriend (keep stored) of %s for day %d", topUser.TopUserIDOutPlace, ownerID, sessDay)
		return nil
	}

	// not in significant places -> check "Best Friend"
	if topUser.TopUserIDOutPlace == friendID {
		// update top user
		topUser.TopUserDurationOutPlace = topUser.TopUserDurationOutPlace + dur
		// put in db
		err = scylladb.PutTopUser(s.dbSession, topUser)
		if err != nil {
			return fmt.Errorf("Failed PutTopUser for %s %s %d: %+v", friendID, ownerID, sessDay, err)
		}
		log.Printf("%s is BestFriend (keep stored but update duration) of %s for day %d", friendID, ownerID, sessDay)
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
		if err != nil {
			return fmt.Errorf("Failed PutTopUser for %s %s %d: %+v", friendID, ownerID, sessDay, err)
		}
		// Best Friend -> ok
		log.Printf("%s is BestFriend (diff) of %s for day %d", friendID, ownerID, sessDay)
		return nil
	}

	// keep top user
	log.Printf("%s is BestFriend (keep stored) of %s for day %d", topUser.TopUserIDOutPlace, ownerID, sessDay)
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
