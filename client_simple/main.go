package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/SUN-XIN/know-your-friends/types"
)

const (
	address = "localhost:8081"
	userID  = "testuser2"
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

	req := types.UserFriendsRequest{
		UserID: userID,
	}
	resp, err := friendsClient.KnowFriends(ctx, &req)
	if err != nil {
		log.Printf("Failed KnowFriends: %v", err)
		return
	}

	log.Printf("%s ->\n\tBestFriend %s\n\tMostSeen %s\n\tCrush %v\n\tMutualLove %s",
		userID,
		resp.BestFriend,
		resp.MostSeen,
		resp.Crush,
		resp.MutualLove)
}
