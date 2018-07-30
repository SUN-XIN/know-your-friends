package main

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/SUN-XIN/know-your-friends/types"
)

const (
	address = "localhost:8081"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	friendsClient := types.NewFriendsClient(conn)

	ctx := context.Background()

	CheckWorkflowCrush(ctx, friendsClient)
	//CheckInitData(ctx, friendsClient)
	//CheckWorkflowBestFriends(ctx, friendsClient)
	//CheckWorkflowMostSeen(ctx, friendsClient)
}

func CheckWorkflowCrush(ctx context.Context, friendsClient types.FriendsClient) {
	var r *types.SessionReply
	// in place (home) + at night + but duration is less than 6h -> only MostSeen
	// night1
	crushUserID := fmt.Sprintf("testuser%dplaces", time.Now().Unix())
	req := types.SessionRequest{
		UserID1:   crushUserID,
		UserID2:   "testuserFriendno1",
		StartDate: 1532566800, // 26/7/2018 03:00:00
		EndDate:   1532574000, // 26/7/2018 05:00:00
		Latitude:  48.823305,
		Longitude: 2.361281,
	}
	_, err := friendsClient.KnowFriends(ctx, &req)
	if err != nil {
		log.Fatalf("Failed KnowFriends: %v", err)
	}

	// night2
	req = types.SessionRequest{
		UserID1:   crushUserID,
		UserID2:   "testuserFriendno1",
		StartDate: 1532653200, // 27/7/2018 03:00:00
		EndDate:   1532660400, // 27/7/2018 03:00:00
		Latitude:  48.823305,
		Longitude: 2.361281,
	}
	_, err = friendsClient.KnowFriends(ctx, &req)
	if err != nil {
		log.Fatalf("Failed KnowFriends: %v", err)
	}

	// night3
	req = types.SessionRequest{
		UserID1:   crushUserID,
		UserID2:   "testuserFriendno1",
		StartDate: 1532739600, // 28/7/2018 03:00:00
		EndDate:   1532746800, // 28/7/2018 05:00:00
		Latitude:  48.823305,
		Longitude: 2.361281,
	}
	r, err = friendsClient.KnowFriends(ctx, &req)
	if err != nil {
		log.Fatalf("Failed KnowFriends: %v", err)
	}

	if r.BestFriend == "" &&
		r.MostSeen == req.UserID2 &&
		len(r.Crush) == 0 {
		log.Printf("Test1 ok")
	} else {
		log.Printf("Test1 not ok: %+v", *r)
	}
	time.Sleep(2 * time.Second)

	// in place (home) + duration 6h30, but not at night-> MostSeen
	// night1
	req = types.SessionRequest{
		UserID1:   crushUserID,
		UserID2:   "testuserFriendno2",
		StartDate: 1532595600, // 26/7/2018 11:00:00
		EndDate:   1532619000, // 26/7/2018 17:30:00
		Latitude:  48.823305,
		Longitude: 2.361281,
	}
	_, err = friendsClient.KnowFriends(ctx, &req)
	if err != nil {
		log.Fatalf("Failed KnowFriends: %v", err)
	}

	// night2
	req = types.SessionRequest{
		UserID1:   crushUserID,
		UserID2:   "testuserFriendno2",
		StartDate: 1532682000, // 27/7/2018 11:00:00
		EndDate:   1532705400, // 27/7/2018 17:30:00
		Latitude:  48.823305,
		Longitude: 2.361281,
	}
	_, err = friendsClient.KnowFriends(ctx, &req)
	if err != nil {
		log.Fatalf("Failed KnowFriends: %v", err)
	}

	// night3
	req = types.SessionRequest{
		UserID1:   crushUserID,
		UserID2:   "testuserFriendno2",
		StartDate: 1532768400, // 28/7/2018 11:00:00
		EndDate:   1532791800, // 28/7/2018 17:30:00
		Latitude:  48.823305,
		Longitude: 2.361281,
	}
	r, err = friendsClient.KnowFriends(ctx, &req)
	if err != nil {
		log.Fatalf("Failed KnowFriends: %v", err)
	}

	if r.BestFriend == "" &&
		r.MostSeen == req.UserID2 &&
		len(r.Crush) == 0 {
		log.Printf("Test2 ok")
	} else {
		log.Printf("Test2 not ok: %+v", *r)
	}
	time.Sleep(2 * time.Second)

	// in place (home) + at night + duration 6h30 -> MostSeen + Crush
	// night1
	req = types.SessionRequest{
		UserID1:   crushUserID,
		UserID2:   "testuserFriend",
		StartDate: 1532566800, // 26/7/2018 3:00:00
		EndDate:   1532590200, // 26/7/2018 9:30:00
		Latitude:  48.823305,
		Longitude: 2.361281,
	}
	_, err = friendsClient.KnowFriends(ctx, &req)
	if err != nil {
		log.Fatalf("Failed KnowFriends: %v", err)
	}

	// night2
	req = types.SessionRequest{
		UserID1:   crushUserID,
		UserID2:   "testuserFriend",
		StartDate: 1532653200, // 27/7/2018 3:00:00
		EndDate:   1532676600, // 27/7/2018 9:30:00
		Latitude:  48.823305,
		Longitude: 2.361281,
	}
	_, err = friendsClient.KnowFriends(ctx, &req)
	if err != nil {
		log.Fatalf("Failed KnowFriends: %v", err)
	}

	// night3
	req = types.SessionRequest{
		UserID1:   crushUserID,
		UserID2:   "testuserFriend",
		StartDate: 1532739600, // 28/7/2018 3:00:00
		EndDate:   1532763500, // 28/7/2018 9:38:20
		Latitude:  48.823305,
		Longitude: 2.361281,
	}
	r, err = friendsClient.KnowFriends(ctx, &req)
	if err != nil {
		log.Fatalf("Failed KnowFriends: %v", err)
	}

	if r.BestFriend == "" &&
		r.MostSeen == req.UserID2 &&
		(len(r.Crush) == 1 && r.Crush[0] == req.UserID2) {
		log.Printf("Test3 ok")
	} else {
		log.Printf("Test3 not ok: %+v", *r)
	}
	time.Sleep(2 * time.Second)

	// another crush
	// night1
	req = types.SessionRequest{
		UserID1:   crushUserID,
		UserID2:   "testuserFriendBis",
		StartDate: 1532566800, // 26/7/2018 3:00:00
		EndDate:   1532590200, // 26/7/2018 9:30:00
		Latitude:  48.823305,
		Longitude: 2.361281,
	}
	_, err = friendsClient.KnowFriends(ctx, &req)
	if err != nil {
		log.Fatalf("Failed KnowFriends: %v", err)
	}

	// night2
	req = types.SessionRequest{
		UserID1:   crushUserID,
		UserID2:   "testuserFriendBis",
		StartDate: 1532653200, // 27/7/2018 3:00:00
		EndDate:   1532676600, // 27/7/2018 9:30:00
		Latitude:  48.823305,
		Longitude: 2.361281,
	}
	_, err = friendsClient.KnowFriends(ctx, &req)
	if err != nil {
		log.Fatalf("Failed KnowFriends: %v", err)
	}

	// night3
	req = types.SessionRequest{
		UserID1:   crushUserID,
		UserID2:   "testuserFriendBis",
		StartDate: 1532739600, // 28/7/2018 3:00:00
		EndDate:   1532763800, // 28/7/2018 9:43:20
		Latitude:  48.823305,
		Longitude: 2.361281,
	}
	r, err = friendsClient.KnowFriends(ctx, &req)
	if err != nil {
		log.Fatalf("Failed KnowFriends: %v", err)
	}

	if r.BestFriend == "" &&
		r.MostSeen == req.UserID2 &&
		len(r.Crush) == 2 {
		log.Printf("Test4 ok: 2 Crush (%v)", r.Crush)
	} else {
		log.Printf("Test4 not ok: %+v", *r)
	}
	time.Sleep(2 * time.Second)
}

func CheckWorkflowMostSeen(ctx context.Context, friendsClient types.FriendsClient) {
	// session 1: int places, not in night -> only MostSeen
	userID := fmt.Sprintf("testuserplaces%d", time.Now().Unix())
	friend1ID := fmt.Sprintf("testfirend%d", time.Now().Unix())
	req := types.SessionRequest{
		UserID1:   userID,
		UserID2:   friend1ID,
		StartDate: 1532419200,
		EndDate:   1532419400,
		Latitude:  48.847016,
		Longitude: 2.355808,
	}
	r, err := friendsClient.KnowFriends(ctx, &req)
	if err != nil {
		log.Fatalf("Failed KnowFriends: %v", err)
	}
	if r.BestFriend == "" &&
		r.MostSeen == friend1ID &&
		len(r.Crush) == 0 {
		log.Printf("Test1 ok")
	} else {
		log.Printf("Test1 not ok: %+v", *r)
	}
	time.Sleep(2 * time.Second)

	// session 2: another friend, more duration + another day
	// in places, not in night -> only MostSeen
	req = types.SessionRequest{
		UserID1:   userID,
		UserID2:   fmt.Sprintf("testfirend%d", time.Now().Unix()),
		StartDate: 1532505600,
		EndDate:   1532506000,
		Latitude:  48.847016,
		Longitude: 2.355808,
	}
	r, err = friendsClient.KnowFriends(ctx, &req)
	if err != nil {
		log.Fatalf("Failed KnowFriends: %v", err)
	}
	if r.BestFriend == "" &&
		r.MostSeen == req.UserID2 &&
		len(r.Crush) == 0 {
		log.Printf("Test2 ok")
	} else {
		log.Printf("Test2 not ok: %+v", *r)
	}
	time.Sleep(2 * time.Second)

	// session 3: friends, more duration + another day
	// in places, not in night -> MostSeen/BestFriend
	req = types.SessionRequest{
		UserID1:   userID,
		UserID2:   friend1ID,
		StartDate: 1532592000,
		EndDate:   1532593000,
		Latitude:  48.847016,
		Longitude: 2.355808,
	}
	r, err = friendsClient.KnowFriends(ctx, &req)
	if err != nil {
		log.Fatalf("Failed KnowFriends: %v", err)
	}
	if r.BestFriend == "" &&
		r.MostSeen == friend1ID &&
		len(r.Crush) == 0 {
		log.Printf("Test3 ok")
	} else {
		log.Printf("Test3 not ok: %+v", *r)
	}
	time.Sleep(2 * time.Second)
}

func CheckWorkflowBestFriends(ctx context.Context, friendsClient types.FriendsClient) {
	// session 1: out places, not in night -> MostSeen/BestFriend
	userID := fmt.Sprintf("testuser%d", time.Now().Unix())
	friend1ID := fmt.Sprintf("testfirend%d", time.Now().Unix())
	req := types.SessionRequest{
		UserID1:   userID,
		UserID2:   friend1ID,
		StartDate: 1532419200,
		EndDate:   1532419400,
		Latitude:  48.847016,
		Longitude: 2.355808,
	}
	r, err := friendsClient.KnowFriends(ctx, &req)
	if err != nil {
		log.Fatalf("Failed KnowFriends: %v", err)
	}
	if r.BestFriend == friend1ID &&
		r.MostSeen == friend1ID &&
		len(r.Crush) == 0 {
		log.Printf("Test1 ok")
	} else {
		log.Printf("Test1 not ok: %+v", *r)
	}
	time.Sleep(2 * time.Second)

	// session 2: another friend, more duration + another day
	// out places, not in night -> MostSeen/BestFriend
	req = types.SessionRequest{
		UserID1:   userID,
		UserID2:   fmt.Sprintf("testfirend%d", time.Now().Unix()),
		StartDate: 1532505600,
		EndDate:   1532506000,
		Latitude:  48.847016,
		Longitude: 2.355808,
	}
	r, err = friendsClient.KnowFriends(ctx, &req)
	if err != nil {
		log.Fatalf("Failed KnowFriends: %v", err)
	}
	if r.BestFriend == req.UserID2 &&
		r.MostSeen == req.UserID2 &&
		len(r.Crush) == 0 {
		log.Printf("Test2 ok")
	} else {
		log.Printf("Test2 not ok: %+v", *r)
	}
	time.Sleep(2 * time.Second)

	// session 3: friends, more duration + another day
	// out places, not in night -> MostSeen/BestFriend
	req = types.SessionRequest{
		UserID1:   userID,
		UserID2:   friend1ID,
		StartDate: 1532592000,
		EndDate:   1532593000,
		Latitude:  48.847016,
		Longitude: 2.355808,
	}
	r, err = friendsClient.KnowFriends(ctx, &req)
	if err != nil {
		log.Fatalf("Failed KnowFriends: %v", err)
	}
	if r.BestFriend == friend1ID &&
		r.MostSeen == friend1ID &&
		len(r.Crush) == 0 {
		log.Printf("Test3 ok")
	} else {
		log.Printf("Test3 not ok: %+v", *r)
	}
	time.Sleep(2 * time.Second)
}

func CheckInitData(ctx context.Context, friendsClient types.FriendsClient) {
	// out place ->  BestFriend + MostSeen
	req := types.SessionRequest{
		UserID1:   fmt.Sprintf("testuser%d", time.Now().Unix()),
		UserID2:   "testuser2",
		StartDate: 1532822000,
		EndDate:   1532822400,
		Latitude:  48.847016,
		Longitude: 2.355808,
	}
	r, err := friendsClient.KnowFriends(ctx, &req)
	if err != nil {
		log.Fatalf("Failed KnowFriends: %v", err)
	}
	if r.BestFriend == req.UserID2 &&
		r.MostSeen == req.UserID2 &&
		len(r.Crush) == 0 {
		log.Printf("Test1 ok")
	} else {
		log.Printf("Test1 not ok: %+v", *r)
	}
	time.Sleep(2 * time.Second)

	// in place (not home) -> only MostSeen
	req = types.SessionRequest{
		UserID1:   fmt.Sprintf("testuser%dplaces", time.Now().Unix()),
		UserID2:   "testuserFriend",
		StartDate: 1532822000,
		EndDate:   1532822400,
		Latitude:  48.847016,
		Longitude: 2.355808,
	}
	r, err = friendsClient.KnowFriends(ctx, &req)
	if err != nil {
		log.Fatalf("Failed KnowFriends: %v", err)
	}
	if r.BestFriend == "" &&
		r.MostSeen == req.UserID2 &&
		len(r.Crush) == 0 {
		log.Printf("Test2 ok")
	} else {
		log.Printf("Test2 not ok: %+v", *r)
	}
	time.Sleep(2 * time.Second)
}
