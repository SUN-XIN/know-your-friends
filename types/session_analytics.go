package types

import "fmt"

const (
	SCYLLA_KEY_SEPARATOR = "-"

	SessionIntegrateTableName = "session_integrate"
	SessionTopUserTableName   = "top_user"
	SessionCrushTableName     = "session_crush"
)

//////////////////////////////////////////////////////////////////
//////////////////////   SessionMostSeen  ////////////////////////
//////////////////////////////////////////////////////////////////

// unit of time is in DAY
type SessionIntegrate struct {
	// key
	UserIDOwner  string
	UserIDFriend string
	Day          int64 // timestamp for the beginning of 1 day

	TotalDuration int32 // in second
	IsInSignPlace bool  // if in my significant place -> only for SessionTypeMostSeen

	//Payload string // json of Payload
}

/*
type Payload struct {
	TotalDuration int32 `json:"total_duration"`   // in second
	IsInSignPlace bool  `json:"is_in_sign_place"` // if in my significant place -> only for SessionTypeMostSeen
}
*/

// format: UserIDOwner-UserIDFriend-Day-Type
func (si *SessionIntegrate) ScyllaDBKey() string {
	return fmt.Sprintf("%s%s%d%s%d", si.UserIDOwner, SCYLLA_KEY_SEPARATOR,
		si.Day)
}

//////////////////////////////////////////////////////////////////
/////////////////////////   SessionCrush  ////////////////////////
//////////////////////////////////////////////////////////////////

// session at home and during night, for Crush
type SessionCrush struct {
	UserIDOwner string
	FriendsIDs  []string
	Day         int64 // timestamp for the beginning of 1 day
}

// sort by UserID, and User1 is min, User2 is max
func (sc *SessionCrush) ScyllaDBKey() string {
	return fmt.Sprintf("%s%s%d", sc.UserIDOwner, SCYLLA_KEY_SEPARATOR,
		sc.Day)
}

/*
// all sessions
type SessionDetail struct {
	UserIDOwner  string
	UserIDFriend string

	StartDate int64
	EndDate   int64

	Lat float64
	Lng float64

	IsInSignPlace bool // if in my significant place
}

// format: UserIDOwner-UserIDFriend-IsInSignPlace
func (sd *SessionDetail) ScyllaDBKey() string {
	return fmt.Sprintf("%s%s%s%s%d", sd.UserIDOwner, SCYLLA_KEY_SEPARATOR,
		sd.UserIDFriend, SCYLLA_KEY_SEPARATOR,
		sd.StartDate)
}
*/
