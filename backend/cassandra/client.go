package cassandra

import (
	"fmt"
	"log"
	// "time"

	"github.com/gocql/gocql"
)

func CassandraClient() {
	fmt.Println("Creating Cassandra session")
	cluster := gocql.NewCluster("cassandra-node1:9042")
	// cluster.Authenticator = gocql.PasswordAuthenticator{
	// 	Username: "user",
	// 	Password: "password"
	// }
	cluster.ProtoVersion = 4
	// cluster.Keyspace = "keyspace"
	// cluster.Port = 29042
	// cluster.CQLVersion = "5.0.0"
	// cluster.ConnectTimeout = time.Second * 6
	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.DCAwareRoundRobinPolicy("cassandra-node1"))
	fmt.Println("3")
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	fmt.Println("Connected to Cassandra Nodes:")
	for _, host := range cluster.Hosts {
		fmt.Println("- Node:", host)
	}
	// session, err := cluster.CreateSession()
	// if err != nil {
	// 	panic(err)
	// }
	// defer session.Close()
	// fmt.Println("5")
	// fmt.Println("Connected to Cassandra Nodes:")
	// for _, host := range cluster.Hosts {
	// 	fmt.Println("- Node:", host)
	// }

	// Print which node is being used for a query
	// 	var clusterID string
	// 	iter := session.Query("SELECT release_version FROM system.local").Iter()
	// 	for iter.Scan(&clusterID) {
	// 		fmt.Println("Running query on:", clusterID)
	// 	}
	// 	if err := iter.Close(); err != nil {
	// 		log.Fatal("Error closing iterator:", err)
	// 	}
}
