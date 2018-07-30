package types

import (
	"fmt"
	"time"

	"github.com/SUN-XIN/know-your-friends/helper"
)

type TopUserType int

type TopUser struct {
	OwnerID string
	Day     int64

	// Best Friend: the person you see the most outside of your significant place
	TopUserIDOutPlace       string
	TopUserDurationOutPlace int32

	// Most seen: the person you see the most
	TopUserID       string
	TopUserDuration int32

	// Crush
	CrushFriendIDs []string
}

// result of MostSeen/BestFriend
func NewTopUserByDuration(ownerID, friendID string, start, end int64, isInPlace bool) *TopUser {
	dur := int32(start - end)
	tu := &TopUser{
		OwnerID:         ownerID,
		Day:             helper.GetBeginningOfDay(time.Unix(end, 0)),
		TopUserID:       friendID,
		TopUserDuration: dur,
	}

	if !isInPlace {
		tu.TopUserIDOutPlace = friendID
		tu.TopUserDurationOutPlace = dur
	}

	return tu
}

func (tu *TopUser) ScyllaDBKey() string {
	return fmt.Sprintf("%s%s%d%s%d", tu.OwnerID, SCYLLA_KEY_SEPARATOR,
		tu.Day)
}
