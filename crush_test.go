package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/SUN-XIN/know-your-friends/types"
)

// 3 nights, but less than 6h
func TestCheckCrush1(t *testing.T) {
	s := NewServer()

	dNow := time.Now().Unix()
	ownerID := fmt.Sprintf("testuserplaces%d", dNow)

	// session1 with friend: less than 6h
	friendID := "testfriendID"
	latHome := 48.823305
	lngHome := 2.361281
	err := s.CheckCrush(ownerID, friendID,
		1532566800, // 26/7/2018 03:00:00
		1532574000, // 26/7/2018 05:00:00
		latHome, lngHome)
	if err != nil {
		t.Errorf("Failed CheckCrush night1: %+v", err)
		return
	}

	// session2 with friend: less than 6h
	err = s.CheckCrush(ownerID, friendID,
		1532653200, // 27/7/2018 03:00:00
		1532660400, // 27/7/2018 03:00:00
		latHome, lngHome)
	if err != nil {
		t.Errorf("Failed CheckCrush night2: %+v", err)
		return
	}

	// session3 with friend: less than 6h
	err = s.CheckCrush(ownerID, friendID,
		1532739600, // 28/7/2018 03:00:00
		1532746800, // 28/7/2018 05:00:00
		latHome, lngHome)
	if err != nil {
		t.Errorf("Failed CheckCrush night3: %+v", err)
		return
	}

	resp := types.UserFriendsReply{}
	err = s.GetCrush(ownerID, dNow, &resp)
	if err != nil {
		t.Errorf("Failed GetCrush: %+v", err)
		return
	}

	if len(resp.Crush) != 0 {
		t.Errorf("Expect Crush empty, but get %+v", resp.Crush)
		return
	}
}

// 3 nights + 6h30 + not in home
// 2 friends
func TestCheckCrush2(t *testing.T) {
	s := NewServer()

	dNow := time.Now().Unix()
	ownerID := fmt.Sprintf("testuserplaces%d", dNow)

	// session1
	friendID := "testfriendID"
	latSchool := 48.847016
	lngSchool := 2.355808
	err := s.CheckCrush(ownerID, friendID,
		1532566800, // 26/7/2018 3:00:00
		1532590200, // 26/7/2018 9:30:00
		latSchool, lngSchool)
	if err != nil {
		t.Errorf("Failed CheckCrush night1: %+v", err)
		return
	}

	// session2
	err = s.CheckCrush(ownerID, friendID,
		1532653200, // 27/7/2018 3:00:00
		1532676600, // 27/7/2018 9:30:00
		latSchool, lngSchool)
	if err != nil {
		t.Errorf("Failed CheckCrush night2: %+v", err)
		return
	}

	// session3
	err = s.CheckCrush(ownerID, friendID,
		1532739600, // 28/7/2018 3:00:00
		1532763500, // 28/7/2018 9:38:20
		latSchool, lngSchool)
	if err != nil {
		t.Errorf("Failed CheckCrush night3: %+v", err)
		return
	}

	resp := types.UserFriendsReply{}
	err = s.GetCrush(ownerID, dNow, &resp)
	if err != nil {
		t.Errorf("Failed GetCrush: %+v", err)
		return
	}

	if len(resp.Crush) != 0 {
		t.Errorf("Expect Crush empty, but get %+v", resp.Crush)
		return
	}
}
