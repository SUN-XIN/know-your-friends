package scylladb

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocql/gocql"

	"github.com/SUN-XIN/know-your-friends/types"
)

func GetTopUser(session *gocql.Session, tu *types.TopUser) error {
	iter := session.Query(fmt.Sprintf(`
	SELECT owner_id, day, top_user_id_out_place, top_user_duration_out_place, to_user_id, top_user_duration, crush_friend_ids 
	FROM %s 
	WHERE owner_id=? AND day=?;`, types.SessionTopUserTableName),
		tu.OwnerID, tu.Day).Iter()
	found := false
	for iter.Scan(&tu.OwnerID,
		&tu.Day,
		&tu.TopUserIDOutPlace,
		&tu.TopUserDurationOutPlace,
		&tu.TopUserID,
		&tu.TopUserDuration,
		&tu.CrushFriendIDs) {
		found = true
		log.Printf("Get ok: %+v", *tu)
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

func PutTopUser(session *gocql.Session, tu *types.TopUser) error {
	if len(tu.CrushFriendIDs) > 0 { // TODO
		return fmt.Errorf("Put with CrushFriendIDs is not supported yet")
	}

	return session.Query(fmt.Sprintf(`
	INSERT INTO %s (owner_id, day, top_user_id_out_place, top_user_duration_out_place, to_user_id, top_user_duration) 
	VALUES (?, ?, ?, ?, ?, ?)`, types.SessionTopUserTableName),
		tu.OwnerID,
		tu.Day,
		tu.TopUserIDOutPlace,
		tu.TopUserDurationOutPlace,
		tu.TopUserID,
		tu.TopUserDuration).Exec()
}

func UpdateTopUserCrush(session *gocql.Session, ownerID, friendID string, day int64) (*types.TopUser, error) {
	tu := &types.TopUser{
		OwnerID: ownerID,
		Day:     day,
	}

	err := GetTopUser(session, tu)
	if err != nil {
		return nil, err
	}

	if friendID != "" {
		alreadyIn := false
		for _, fID := range tu.CrushFriendIDs {
			if fID == friendID {
				alreadyIn = true
				break
			}
		}

		if !alreadyIn {
			tu.CrushFriendIDs = append(tu.CrushFriendIDs, friendID)
		}
	}

	vals := make([]string, 0, len(tu.CrushFriendIDs))
	for _, f := range tu.CrushFriendIDs {
		vals = append(vals, fmt.Sprintf("'%s'", f))
	}

	session.Query(fmt.Sprintf(`
	UPDATE %s SET crush_friend_ids = {%s}
	WHERE owner_id=? AND day=?`, types.SessionTopUserTableName, strings.Join(vals, ",")),
		tu.OwnerID,
		tu.Day).Exec()
	return tu, err
}
