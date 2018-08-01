package types

import "fmt"

const (
	SCYLLA_KEY_SEPARATOR = "-"

	SessionIntegrateTableName = "session_integrate"
	SessionTopUserTableName   = "top_user"
	SessionCrushTableName     = "session_crush"
	SessionDetailTableName    = "session_detail"
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
	return fmt.Sprintf("%s%s%s%s%d", si.UserIDOwner, SCYLLA_KEY_SEPARATOR,
		si.UserIDFriend, SCYLLA_KEY_SEPARATOR,
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

// all sessions
type SessionDetail struct {
	UserID1 string
	UserID2 string

	StartDate int64
	EndDate   int64

	Lat float64
	Lng float64
}

// Check if empty for each field
// UserID1 should always be greater than UserID2
func (sd *SessionDetail) Validate() error {
	if sd.UserID1 == sd.UserID2 {
		return fmt.Errorf("UserIDs must not be diffenrent")
	}

	if sd.UserID1 == "" || sd.UserID2 == "" {
		return fmt.Errorf("UserID must not be empty")
	}

	if sd.StartDate <= 0 || sd.EndDate <= 0 ||
		sd.EndDate <= sd.StartDate {
		return fmt.Errorf("Please check startDate/EndDate")
	}

	if sd.Lat < -90.0 || sd.Lat > 90.0 ||
		sd.Lng < -180.0 || sd.Lng > 180.0 {
		return fmt.Errorf("Please check Lat/Lng")
	}

	if sd.UserID1 < sd.UserID2 {
		sd.UserID1, sd.UserID2 = sd.UserID2, sd.UserID1
	}

	return nil
}

// format: UserIDOwner-UserIDFriend-StartDate
func (sd *SessionDetail) ScyllaDBKey() string {
	return fmt.Sprintf("%s%s%s%s%d", sd.UserID1, SCYLLA_KEY_SEPARATOR,
		sd.UserID2, SCYLLA_KEY_SEPARATOR,
		sd.StartDate)
}
