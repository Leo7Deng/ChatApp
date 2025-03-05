package cassandra

import (
	"fmt"
	"log"
	// "time"

	"github.com/gocql/gocql"
)

var (
	cluster *gocql.ClusterConfig
	session *gocql.Session
)

func CassandraSession() *gocql.Session {
	fmt.Println("Creating Cassandra session")
	cluster = gocql.NewCluster("cassandra-node1:9042", "cassandra-node2:9042", "cassandra-node3:9042")
	var err error
	cluster.Consistency = gocql.Quorum
	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.DCAwareRoundRobinPolicy("cassandra-node1"))
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Cassandra Nodes:")
	for _, host := range cluster.Hosts {
		fmt.Println("- Node:", host)
	}

	// Print which node is being used for a query
	var clusterID string
	iter := session.Query("SELECT release_version FROM system.local").Iter()
	for iter.Scan(&clusterID) {
		fmt.Println("Running query on node ", clusterID)
	}
	if err := iter.Close(); err != nil {
		log.Fatal("Error closing iterator:", err)
	}

	return session
}

func Close() {
	if session != nil {
		session.Close()
		fmt.Println("Closing Cassandra session")
	}
}