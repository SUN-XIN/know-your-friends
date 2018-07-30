package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/SUN-XIN/know-your-friends/types"
)

const (
	DEFAULT_USERS_SIGNIFICANT_PLACE_SIZE = 1000
	DEFAULT_USERS_TOP_SIZE               = 1000

	DEFAULT_USERS_NUM = 10000
)

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	serv := NewServer()
	types.RegisterFriendsServer(s, serv)
	s.Serve(lis)
}
