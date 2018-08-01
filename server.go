package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/gocql/gocql"
	lru "github.com/hashicorp/golang-lru"
	"golang.org/x/net/context"

	"github.com/SUN-XIN/know-your-friends/helper"
	"github.com/SUN-XIN/know-your-friends/scylladb"
	"github.com/SUN-XIN/know-your-friends/types"
)

const (
	DEFAULT_FRIENDS_PER_DAY = 100
	KAFKA_BROKER            = "localhost:9092"
	KAFKA_TOPIC             = "test_zenly"
	KAFKA_GROUP             = "g1"
)

type server struct {
	dbSession   *gocql.Session
	configKafka *cluster.Config

	sessionReceived *lru.Cache // avoid of double session
	cacheUserPlaces *lru.Cache // cache for users' significant places
	cacheUserTop    *lru.Cache // cache of result

	mutex *sync.Mutex
}

func NewServer() *server {
	// size never < 0, we can ignore error
	cacheSignPlace, _ := lru.New(DEFAULT_USERS_SIGNIFICANT_PLACE_SIZE)
	cacheTopUser, _ := lru.New(DEFAULT_USERS_TOP_SIZE)
	sessionReceived, _ := lru.New(DEFAULT_USERS_TOP_SIZE)
	dbSess, err := scylladb.CreateConnect()
	if err != nil {
		panic(err.Error())
	}

	config := cluster.NewConfig()
	config.Group.Return.Notifications = true
	config.Consumer.Return.Errors = true
	config.Version = sarama.MaxVersion

	serv := &server{
		dbSession: dbSess,

		configKafka: config,

		sessionReceived: sessionReceived,
		cacheUserPlaces: cacheSignPlace,
		cacheUserTop:    cacheTopUser,

		mutex: &sync.Mutex{},
	}

	go func() {
		brokers := []string{KAFKA_BROKER}
		topics := []string{KAFKA_TOPIC}
		consumer, err := cluster.NewConsumer(brokers, KAFKA_GROUP, topics, config)
		if err != nil {
			panic(err)
		}
		defer consumer.Close()

		go func(c *cluster.Consumer) {
			errors := c.Errors()
			notif := c.Notifications()
			for {
				select {
				case err := <-errors:
					log.Printf("Errors chan: %+v", err)
				case not := <-notif:
					log.Printf("Notifications chan: %+v", not)
				}
			}
		}(consumer)

		for msg := range consumer.Messages() {
			log.Printf("%s/%d/%d\t%s", msg.Topic, msg.Partition, msg.Offset, msg.Value)
			consumer.MarkOffset(msg, "")

			var req types.SessionDetail
			err = json.Unmarshal(msg.Value, &req)
			if err != nil {
				log.Printf("Failed Unmarshal: %+v", err)
				continue
			}

			err = serv.ProcessSession(&req)
			if err != nil {
				log.Printf("Failed ProcessSession for (%+v): %+v", req, err)
				continue
			}

			log.Printf("%v process ok", req)
		}
	}()

	return serv
}

func (serv *server) CheckDoubleSession(sd *types.SessionDetail) error {
	serv.mutex.Lock()
	defer serv.mutex.Unlock()

	err := sd.Validate()
	if err != nil {
		return fmt.Errorf("Bad request: +v", err)
	}

	return scylladb.CreateSessionDetail(serv.dbSession, sd)
}

func (serv *server) KnowFriends(ctx context.Context, sess *types.UserFriendsRequest) (*types.UserFriendsReply, error) {
	log.Printf("receive user %s", sess.UserID)

	resp := types.UserFriendsReply{}
	today := helper.GetBeginningOfDay(time.Now().UTC())
	err := serv.GetBestFriendAndMostSeen(sess.UserID, today, &resp)
	if err != nil {
		log.Printf("Failed GetBestFriendAndMostSeen: %+v", err)
		return nil, fmt.Errorf("Failed GetBestFriendAndMostSeen: %+v", err)
	}
	log.Printf("user %s GetBestFriendAndMostSeen ok", sess.UserID)

	// Mutual Love
	respFriend := types.UserFriendsReply{}
	err = serv.GetBestFriendAndMostSeen(resp.MostSeen, today, &respFriend)
	if err != nil {
		log.Printf("Failed GetBestFriendAndMostSeen: %+v", err)
		return nil, fmt.Errorf("Failed GetBestFriendAndMostSeen for friend: %+v", err)
	}
	if respFriend.MostSeen == sess.UserID {
		resp.MutualLove = resp.MostSeen
	}
	log.Printf("user %s MutualLove ok", sess.UserID)

	err = serv.GetCrush(sess.UserID, today, &resp)
	if err != nil {
		log.Printf("Failed GetCrush: %+v", err)
		return nil, fmt.Errorf("Failed GetCrush: %+v", err)
	}
	log.Printf("user %s GetCrush ok", sess.UserID)

	return &resp, nil
}

func (serv *server) ProcessSession(sess *types.SessionDetail) error {
	// session too old
	if sess.EndDate < time.Now().Add(-24*time.Hour*helper.DEFAULT_ROLLING_DAYS).Unix() {
		return nil
	}

	// double session
	err := serv.CheckDoubleSession(sess)
	switch err {
	case nil:
	case scylladb.ErrAlreadyExist:
		log.Printf("Session already received, nothing to do")
		return fmt.Errorf("Session already received, nothing to do")
	default:
		log.Printf("Failed CheckDoubleSession: %+v", err)
		return fmt.Errorf("Failed CheckDoubleSession: %+v", err)
	}

	var goErr error
	wg := sync.WaitGroup{}
	wg.Add(4)
	go func(e *error) {
		defer wg.Done()

		err := serv.CheckBestFriendAndMostSeen(sess.UserID1, sess.UserID2,
			sess.Lat, sess.Lng,
			sess.StartDate, sess.EndDate)
		if err != nil {
			log.Printf("Failed CheckBestFriendAndMostSeen: %+v", err)
			e = &err
			return
		}
	}(&goErr)

	go func(e *error) {
		defer wg.Done()

		err := serv.CheckBestFriendAndMostSeen(sess.UserID2, sess.UserID1,
			sess.Lat, sess.Lng,
			sess.StartDate, sess.EndDate)
		if err != nil {
			log.Printf("Failed CheckBestFriendAndMostSeen: %+v", err)
			e = &err
			return
		}
	}(&goErr)

	go func(e *error) {
		defer wg.Done()

		err := serv.CheckCrush(sess.UserID1, sess.UserID2,
			sess.StartDate, sess.EndDate,
			sess.Lat, sess.Lng)
		if err != nil {
			log.Printf("Failed CheckCrush: %+v", err)
			e = &err
		}
	}(&goErr)

	go func(e *error) {
		defer wg.Done()

		err := serv.CheckCrush(sess.UserID2, sess.UserID1,
			sess.StartDate, sess.EndDate,
			sess.Lat, sess.Lng)
		if err != nil {
			log.Printf("Failed CheckCrush: %+v", err)
			e = &err
		}
	}(&goErr)

	wg.Wait()
	return goErr
}

/*
func (s *server) CheckMutualLove(sess *types.SessionRequest,
	resp *types.SessionReply) error {
	if resp.User1MostSeen != sess.UserID2 &&
		resp.User2MostSeen != sess.UserID1 {
		resp.User1MutualLove = sess.UserID2
		resp.User2MutualLove = sess.UserID1
	}
	return nil
}
*/
