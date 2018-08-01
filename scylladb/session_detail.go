package scylladb

import (
	"fmt"
	"log"

	"github.com/SUN-XIN/know-your-friends/types"
	"github.com/gocql/gocql"
)

func GetSessionDetail(session *gocql.Session, sd *types.SessionDetail) error {
	iter := session.Query(fmt.Sprintf(`
	SELECT user_id_1, user_id_2, start_date, end_date, lat, lng
	FROM %s 
	WHERE user_id_1=? AND user_id_2=? AND start_date=?;`, types.SessionDetailTableName),
		sd.UserID1, sd.UserID2, sd.StartDate).Iter()
	found := false
	for iter.Scan(&sd.UserID1,
		&sd.UserID2,
		&sd.StartDate,
		&sd.EndDate,
		&sd.Lat,
		&sd.Lng) {
		found = true
		log.Printf("Get ok: %+v", *sd)
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

func PutSessionDetail(session *gocql.Session, sd *types.SessionDetail) error {
	return session.Query(fmt.Sprintf(`
	INSERT INTO %s (user_id_1,user_id_2,start_date,end_date,lat,lng) 
	VALUES (?, ?, ?, ?, ?, ?)`, types.SessionDetailTableName),
		sd.UserID1,
		sd.UserID2,
		sd.StartDate,
		sd.EndDate,
		sd.Lat,
		sd.Lng).Exec()
}

func CreateSessionDetail(session *gocql.Session, sd *types.SessionDetail) error {
	err := GetSessionDetail(session, sd)
	switch err {
	case nil:
		return ErrAlreadyExist
	case ErrNotFound:
		// ok
	default:
		return err
	}

	log.Printf("Create -> get ok, not found")
	return PutSessionDetail(session, sd)
}
