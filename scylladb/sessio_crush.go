package scylladb

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocql/gocql"

	"github.com/SUN-XIN/know-your-friends/helper"
	"github.com/SUN-XIN/know-your-friends/types"
)

func GetSessionCrush(session *gocql.Session, sr *types.SessionCrush) error {
	iter := session.Query(fmt.Sprintf(`
	SELECT user_id_owner,user_ids_friend,day 
	FROM %s 
	WHERE user_id_owner=? AND day=?;`, types.SessionCrushTableName),
		sr.UserIDOwner, sr.Day).Iter()
	found := false

	for iter.Scan(&sr.UserIDOwner,
		&sr.FriendsIDs,
		&sr.Day) {
		found = true
		log.Printf("Get ok: %+v", *sr)
	}
	err := iter.Close()
	if err != nil {
		return err
	}

	if !found {
		return ErrNotFound
	}

	return nil
}

func GetAndUpdateSessionCrush(session *gocql.Session, ownerID, friendID string, day int64) (*types.SessionCrush, error) {
	iter := session.Query(fmt.Sprintf(`
	SELECT user_id_owner,user_ids_friend,day 
	FROM %s 
	WHERE user_id_owner=? AND day=?;`, types.SessionCrushTableName),
		ownerID, day).Iter()
	found := false

	sr := &types.SessionCrush{}
	for iter.Scan(&sr.UserIDOwner,
		&sr.FriendsIDs,
		&sr.Day) {
		found = true
		log.Printf("Get ok: %+v", *sr)
	}
	err := iter.Close()
	if err != nil {
		return nil, err
	}

	if !found { // create
		sr.UserIDOwner = ownerID
		sr.FriendsIDs = []string{friendID}
		sr.Day = day

		return sr, PutSessionCrush(session, sr)
	}

	// update ?
	found = false
	for _, fID := range sr.FriendsIDs {
		if fID == friendID {
			found = true
			break
		}
	}
	if found { // session already existed ?
		log.Printf("SHOULD NEVER HAPPEN: session crush alrady existed")
		return sr, nil
	}

	sr.FriendsIDs = append(sr.FriendsIDs, friendID)
	return sr, PutSessionCrush(session, sr)
}

func PutSessionCrush(session *gocql.Session, sr *types.SessionCrush) error {
	lFriends := len(sr.FriendsIDs)
	vals := make([]string, 0, lFriends)

	for i := 0; i < lFriends; i++ {
		vals = append(vals, fmt.Sprintf("'%s'", sr.FriendsIDs[i]))
	}

	return session.Query(fmt.Sprintf(`
		INSERT INTO %s (user_id_owner,day,user_ids_friend) 
		VALUES (?, ?, {%s})`, types.SessionCrushTableName, strings.Join(vals, ",")),
		sr.UserIDOwner, sr.Day).Exec()
}

func CountNights(session *gocql.Session, si *types.SessionCrush) ([]string, error) {
	var err error
	// how many SessionHomeNight in the last 7 days
	nb := 0
	days := helper.GetLastDays(time.Now())
	mapCrush := make(map[string]int, 100)
	var ok bool
	for _, d := range days {
		pastSession := types.SessionCrush{
			UserIDOwner: si.UserIDOwner,
			Day:         d,
		}

		err = GetSessionCrush(session, &pastSession)
		switch err {
		case nil:
			for _, fID := range pastSession.FriendsIDs {
				nb, ok = mapCrush[fID]
				if ok {
					nb++
					mapCrush[fID] = nb
				} else {
					mapCrush[fID] = 1
				}
			}
		case ErrNotFound:
		default:
			return nil, err
		}
	}

	res := make([]string, 0, len(mapCrush))
	for fID, count := range mapCrush {
		if count >= helper.CRUSH_MIN_NIGHTS {
			log.Printf("%s is Crush of %s", fID, si.UserIDOwner)
			res = append(res, fID)
		}
	}
	log.Printf("Total %d Crush of %s", len(res), si.UserIDOwner)

	return res, nil
}
