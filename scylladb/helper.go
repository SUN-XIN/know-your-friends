package scylladb

import (
	"fmt"

	"github.com/gocql/gocql"
)

const (
	CLUSTER  = "172.17.0.1:9042"
	KEYSPASE = "friends"
)

var (
	ErrNotFound = fmt.Errorf("EntityNotFound")
)

func CreateConnect() (*gocql.Session, error) {
	cluster := gocql.NewCluster(CLUSTER)
	cluster.Keyspace = KEYSPASE
	cluster.Consistency = gocql.Quorum
	return cluster.CreateSession()
}
