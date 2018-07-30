package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
	lru "github.com/hashicorp/golang-lru"
	"golang.org/x/net/context"

	"github.com/SUN-XIN/know-your-friends/helper"
	"github.com/SUN-XIN/know-your-friends/scylladb"
	"github.com/SUN-XIN/know-your-friends/types"
)

type server struct {
	dbSession *gocql.Session

	cacheUserPlaces *lru.Cache // cache for users' significant places
	cacheUserTop    *lru.Cache // cache of result
}

func NewServer() *server {
	// size never < 0, we can ignore error
	cacheSignPlace, _ := lru.New(DEFAULT_USERS_SIGNIFICANT_PLACE_SIZE)
	cacheTopUser, _ := lru.New(DEFAULT_USERS_TOP_SIZE)

	dbSess, err := scylladb.CreateConnect()
	if err != nil {
		panic(err.Error())
	}
	return &server{
		dbSession: dbSess,

		cacheUserPlaces: cacheSignPlace,
		cacheUserTop:    cacheTopUser,
	}
}

func (serv *server) KnowFriends(ctx context.Context, sess *types.SessionRequest) (*types.SessionReply, error) {
	resp := types.SessionReply{}

	// session too old
	if sess.EndDate < time.Now().Add(-24*time.Hour*helper.DEFAULT_ROLLING_DAYS).Unix() {
		return &resp, nil
	}

	err := serv.CheckBestFriendAndMostSeen(sess.UserID1, sess.UserID2,
		sess.Latitude, sess.Longitude,
		sess.StartDate, sess.EndDate,
		&resp)
	if err != nil {
		log.Printf("Failed CheckBestFriendAndMostSeen: %+v", err)
		return nil, fmt.Errorf("Failed CheckBestFriendAndMostSeen: +v", err)
	}

	err = serv.CheckCrush(sess, &resp)
	if err != nil {
		log.Printf("Failed CheckCrush: %+v", err)
		return nil, fmt.Errorf("Failed CheckCrush: +v", err)
	}

	err = serv.CheckMutualLove(sess, &resp)
	if err != nil {
		log.Printf("Failed CheckMutualLove: %+v", err)
		return nil, fmt.Errorf("Failed CheckMutualLove: +v", err)
	}

	/*
		keys := serv.cacheUserTop.Keys()
		for _, k := range keys {
			t, _ := serv.cacheUserTop.Get(k)
			tu := t.(*types.TopUser)
			log.Printf("owner %s: mostSeenID %s dur %d", tu.OwnerID, tu.TopUserID, tu.TopUserDuration)
		}
	*/

	return &resp, nil
}

func (s *server) CheckMutualLove(sess *types.SessionRequest,
	resp *types.SessionReply) error {
	if resp.MostSeen != sess.UserID2 {
		log.Printf("MostSeen is another User")
		return nil
	}
	// friend's top user
	topUserInterf, existed := s.cacheUserTop.Get(sess.UserID2)
	if !existed {
		log.Printf("friend's topUser not found TopUser in cache")
		return nil
	}

	topUserFriend, ok := topUserInterf.(*types.TopUser)
	if !ok {
		return fmt.Errorf("SHOULD NEVER HAPPEN, Failed Conv to TopUser: %+v", topUserInterf)
	}

	log.Printf("TopUserID %s UserID1 %s", topUserFriend.TopUserID, sess.UserID1)
	if topUserFriend.TopUserID == sess.UserID1 {
		resp.MutualLove = sess.UserID2
	}

	return nil
}
