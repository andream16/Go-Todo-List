package cassandramanager

import (
	"github.com/gocql/gocql"
	"errors"
)

func InitCassandraCluster() (gocql.Session, error) {

	var session *gocql.Session

	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "todos"
	cluster.Consistency = gocql.One

	if session == nil || session.Closed() {
		session, err := cluster.CreateSession()
		if( err != nil ){
			return *session, errors.New("Cannot Initialize Cassandra Cluster : " + err.Error())
		}

		return *session, err
	} else {
		return *session, nil
	}

}
