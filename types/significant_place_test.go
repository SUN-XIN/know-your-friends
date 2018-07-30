package types

import (
	"fmt"
	"testing"
	"time"

	"github.com/SUN-XIN/know-your-friends/geo"
)

func TestIsIn(t *testing.T) {
	userID := "testName"

	inHomeLat, inHomeLng := 48.823305, 2.361281
	pl := SignificantPlace{
		UserID: userID,
		ID:     fmt.Sprintf("%d", time.Now().Unix()),
		LLBox: geo.LLBox{
			N: 48.823534,
			E: 2.361528,
			S: 48.823058,
			W: 2.360991,
		},
		Name: PLACE_NAME_HOME,
	}
	if !pl.IsIn(inHomeLat, inHomeLng) {
		t.Errorf("Expect in home")
	}

	inSchoolLat, inSchoolLng := 48.847016, 2.355808
	pl = SignificantPlace{
		UserID: userID,
		ID:     fmt.Sprintf("%d", time.Now().Unix()),
		LLBox: geo.LLBox{
			N: 48.849049,
			E: 2.357621,
			S: 48.844919,
			W: 2.355068,
		},
		Name: PLACE_NAME_SCHOOL,
	}
	if !pl.IsIn(inSchoolLat, inSchoolLng) {
		t.Errorf("Expect in school")
	}

	inWorkLat, inWorkLng := 48.854095, 2.373171
	pl = SignificantPlace{
		UserID: userID,
		ID:     fmt.Sprintf("%d", time.Now().Unix()),
		LLBox: geo.LLBox{
			N: 48.854147,
			E: 2.373185,
			S: 48.853667,
			W: 2.372981,
		},
		Name: PLACE_NAME_WORK,
	}
	if !pl.IsIn(inWorkLat, inWorkLng) {
		t.Errorf("Expect in work")
	}

}
