package scylladb

import (
	"fmt"
	"log"

	"github.com/gocql/gocql"

	"github.com/SUN-XIN/know-your-friends/types"
)

func GetSessionIntegrate(session *gocql.Session, si *types.SessionIntegrate) error {
	iter := session.Query(fmt.Sprintf(`
	SELECT user_id_owner,user_id_friend,day,is_in_sign_place,total_duration 
	FROM %s 
	WHERE user_id_owner=? AND user_id_friend=? AND day=?;`, types.SessionIntegrateTableName),
		si.UserIDOwner, si.UserIDFriend, si.Day).Iter()
	found := false
	for iter.Scan(&si.UserIDOwner,
		&si.UserIDFriend,
		&si.Day,
		&si.IsInSignPlace,
		&si.TotalDuration) {
		found = true
		log.Printf("Get ok: %+v", *si)
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

func PutSessionIntegrate(session *gocql.Session, si *types.SessionIntegrate) error {
	return session.Query(fmt.Sprintf(`
	INSERT INTO %s (user_id_owner,user_id_friend,day,is_in_sign_place,total_duration) 
	VALUES (?, ?, ?, ?, ?)`, types.SessionIntegrateTableName),
		si.UserIDOwner,
		si.UserIDFriend,
		si.Day,
		si.IsInSignPlace,
		si.TotalDuration).Exec()
}

func CreateSessionIntegrate(session *gocql.Session, si *types.SessionIntegrate) error {
	err := GetSessionIntegrate(session, si)
	if err == nil {
		return fmt.Errorf("Already existed")
	}
	return PutSessionIntegrate(session, si)
}

func UpdateSessionIntegrate(session *gocql.Session, si *types.SessionIntegrate, dur int32) error {
	err := GetSessionIntegrate(session, si)
	switch err {
	case nil: // existed -> update
		si.TotalDuration = si.TotalDuration + dur
	case ErrNotFound: // not found -> create
		si.TotalDuration = dur
	default:
		return fmt.Errorf("Failed GetSessionIntegrate: %+v", err)
	}

	err = PutSessionIntegrate(session, si)
	if err != nil {
		return fmt.Errorf("Failed PutSessionIntegrate: %+v", err)
	}

	return nil
}

//////////////////////////////////////////////////////////////////////////////////
/////////////////////////////    SessionDetail       /////////////////////////////
//////////////////////////////////////////////////////////////////////////////////
/*
func GetSessionDetail(sd *types.SessionDetail) error {
	sd.ScyllaDBKey()
	return nil
}

func PutSessionDetail(sd *types.SessionDetail) error {
	sd.ScyllaDBKey()
	return nil
}

func CreateSessionDetail(sd *types.SessionDetail) error {
	// transaction
	err := GetSessionDetail(sd)
	if err == nil {
		return fmt.Errorf("Already existed")
	}
	return PutSessionDetail(sd)
}

func FetchAllSessionDetailOfDay(day int64) ([]*types.SessionDetail, error) {
	// Query with: UserIDOwner + UserIDFriend + StartDate
	return nil, nil
}
*/
