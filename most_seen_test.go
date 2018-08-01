package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/SUN-XIN/know-your-friends/helper"
	"github.com/SUN-XIN/know-your-friends/types"
)

// Test MostSeen: session is in some place
func TestCheckBestFriendAndMostSeen1(t *testing.T) {
	s := NewServer()

	dNow := time.Now().Unix()
	ownerID := fmt.Sprintf("testuserplaces%d", dNow)

	// session1 with friend1: 300s
	friendID1 := "testfriendID1"
	latSchool := 48.847016
	lngSchool := 2.355808
	err := s.CheckBestFriendAndMostSeen(ownerID, friendID1,
		latSchool, lngSchool,
		dNow-1000, dNow-700)
	if err != nil {
		t.Errorf("Failed CheckBestFriendAndMostSeen: %+v", err)
		return
	}

	cacheKey := fmt.Sprintf("%s-%d", ownerID, helper.GetBeginningOfDay(time.Unix(dNow, 0)))
	topUserInterf, ok := s.cacheUserTop.Get(cacheKey)
	if !ok {
		t.Errorf("Not found UserTop in cache")
		return
	}

	topUser, ok := topUserInterf.(*types.TopUser)
	if !ok {
		t.Errorf("Failed conv")
		return
	}

	if topUser.TopUserID != friendID1 {
		t.Errorf("Expect MostSeen %s, but get %s", friendID1, topUser.TopUserID)
	}
	if topUser.TopUserDuration != 300 {
		t.Errorf("Expect Duration MostSeen 300, but get %s", topUser.TopUserDuration)
	}
	if topUser.TopUserIDOutPlace != "" {
		t.Errorf("Expect BestFriend empty, but get %s", topUser.TopUserIDOutPlace)
	}
	if topUser.TopUserDurationOutPlace != 0 {
		t.Errorf("Expect Duration BestFriend 0, but get %s", topUser.TopUserDurationOutPlace)
	}
	t.Logf("session1 with friend1 300s -> done")

	// session2 with friend2: 200s
	friendID2 := "testfriendID2"
	err = s.CheckBestFriendAndMostSeen(ownerID, friendID2,
		latSchool, lngSchool,
		dNow-700, dNow-500)
	if err != nil {
		t.Errorf("Failed CheckBestFriendAndMostSeen: %+v", err)
		return
	}

	topUserInterf, ok = s.cacheUserTop.Get(cacheKey)
	if !ok {
		t.Errorf("Not found UserTop in cache")
		return
	}

	topUser, ok = topUserInterf.(*types.TopUser)
	if !ok {
		t.Errorf("Failed conv")
		return
	}

	if topUser.TopUserID != friendID1 {
		t.Errorf("Expect MostSeen %s, but get %s", friendID1, topUser.TopUserID)
	}
	if topUser.TopUserDuration != 300 {
		t.Errorf("Expect Duration MostSeen 300, but get %s", topUser.TopUserDuration)
	}
	if topUser.TopUserIDOutPlace != "" {
		t.Errorf("Expect BestFriend empty, but get %s", topUser.TopUserIDOutPlace)
	}
	if topUser.TopUserDurationOutPlace != 0 {
		t.Errorf("Expect Duration BestFriend 0, but get %s", topUser.TopUserDurationOutPlace)
	}
	t.Logf("session2 with friend2 200s -> done")

	// session3 with friend2: 200s
	err = s.CheckBestFriendAndMostSeen(ownerID, friendID2,
		latSchool, lngSchool,
		dNow-500, dNow-300)
	if err != nil {
		t.Errorf("Failed CheckBestFriendAndMostSeen: %+v", err)
		return
	}

	topUserInterf, ok = s.cacheUserTop.Get(cacheKey)
	if !ok {
		t.Errorf("Not found UserTop in cache")
		return
	}

	topUser, ok = topUserInterf.(*types.TopUser)
	if !ok {
		t.Errorf("Failed conv")
		return
	}

	if topUser.TopUserID != friendID2 {
		t.Errorf("Expect MostSeen %s, but get %s", friendID1, topUser.TopUserID)
	}
	if topUser.TopUserDuration != 400 {
		t.Errorf("Expect Duration MostSeen 400, but get %d", topUser.TopUserDuration)
	}
	if topUser.TopUserIDOutPlace != "" {
		t.Errorf("Expect BestFriend empty, but get %s", topUser.TopUserIDOutPlace)
	}
	if topUser.TopUserDurationOutPlace != 0 {
		t.Errorf("Expect Duration BestFriend 0, but get %d", topUser.TopUserDurationOutPlace)
	}
	t.Logf("session3 with friend2 200s -> done")
}

// Test MostSeen and BestFriend: session is out of places
func TestCheckBestFriendAndMostSeen2(t *testing.T) {
	s := NewServer()

	dNow := time.Now().Unix()
	ownerID := fmt.Sprintf("testmsbf%d", dNow)

	// session1 with firned1: 300
	friendID1 := "testfriendID1"
	latSchool := 38.847016
	lngSchool := 12.355808
	err := s.CheckBestFriendAndMostSeen(ownerID, friendID1,
		latSchool, lngSchool,
		dNow-1000, dNow-700)
	if err != nil {
		t.Errorf("Failed CheckBestFriendAndMostSeen: %+v", err)
		return
	}

	cacheKey := fmt.Sprintf("%s-%d", ownerID, helper.GetBeginningOfDay(time.Unix(dNow, 0)))
	topUserInterf, ok := s.cacheUserTop.Get(cacheKey)
	if !ok {
		t.Errorf("Not found UserTop in cache")
		return
	}

	topUser, ok := topUserInterf.(*types.TopUser)
	if !ok {
		t.Errorf("Failed conv")
		return
	}

	if topUser.TopUserID != friendID1 {
		t.Errorf("Expect MostSeen %s, but get %s", friendID1, topUser.TopUserID)
	}
	if topUser.TopUserDuration != 300 {
		t.Errorf("Expect Duration MostSeen 300, but get %d", topUser.TopUserDuration)
	}
	if topUser.TopUserIDOutPlace != friendID1 {
		t.Errorf("Expect BestFriend %s, but get %s", friendID1, topUser.TopUserIDOutPlace)
	}
	if topUser.TopUserDurationOutPlace != 300 {
		t.Errorf("Expect Duration BestFriend 300, but get %d", topUser.TopUserDurationOutPlace)
	}
	t.Logf("session1 with friend1 300s -> done")

	// session2 with firned2: 200
	friendID2 := "testfriendID2"
	err = s.CheckBestFriendAndMostSeen(ownerID, friendID2,
		latSchool, lngSchool,
		dNow-700, dNow-500)
	if err != nil {
		t.Errorf("Failed CheckBestFriendAndMostSeen: %+v", err)
		return
	}

	topUserInterf, ok = s.cacheUserTop.Get(cacheKey)
	if !ok {
		t.Errorf("Not found UserTop in cache")
		return
	}

	topUser, ok = topUserInterf.(*types.TopUser)
	if !ok {
		t.Errorf("Failed conv")
		return
	}

	if topUser.TopUserID != friendID1 {
		t.Errorf("Expect MostSeen %s, but get %s", friendID1, topUser.TopUserID)
	}
	if topUser.TopUserDuration != 300 {
		t.Errorf("Expect Duration MostSeen 300, but get %d", topUser.TopUserDuration)
	}
	if topUser.TopUserIDOutPlace != friendID1 {
		t.Errorf("Expect BestFriend %s, but get %s", friendID1, topUser.TopUserIDOutPlace)
	}
	if topUser.TopUserDurationOutPlace != 300 {
		t.Errorf("Expect Duration BestFriend 300, but get %d", topUser.TopUserDurationOutPlace)
	}
	t.Logf("session2 with friend2 200s -> done")

	// session3 with firned1: 200
	err = s.CheckBestFriendAndMostSeen(ownerID, friendID2,
		latSchool, lngSchool,
		dNow-500, dNow-300)
	if err != nil {
		t.Errorf("Failed CheckBestFriendAndMostSeen: %+v", err)
		return
	}

	topUserInterf, ok = s.cacheUserTop.Get(cacheKey)
	if !ok {
		t.Errorf("Not found UserTop in cache")
		return
	}

	topUser, ok = topUserInterf.(*types.TopUser)
	if !ok {
		t.Errorf("Failed conv")
		return
	}

	if topUser.TopUserID != friendID2 {
		t.Errorf("Expect MostSeen %s, but get %s", friendID2, topUser.TopUserID)
	}
	if topUser.TopUserDuration != 400 {
		t.Errorf("Expect Duration MostSeen 400, but get %d", topUser.TopUserDuration)
	}
	if topUser.TopUserIDOutPlace != friendID2 {
		t.Errorf("Expect BestFriend %s, but get %s", friendID2, topUser.TopUserIDOutPlace)
	}
	if topUser.TopUserDurationOutPlace != 400 {
		t.Errorf("Expect Duration BestFriend 400, but get %d", topUser.TopUserDurationOutPlace)
	}
	t.Logf("session3 with friend2 200s -> done")
}
