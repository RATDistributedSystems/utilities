package ratdatabase

import (
	"log"

	"github.com/gocql/gocql"
)

var cassandraConnection *gocql.Session

// InitCassandraConnection is called by a server to initialize the connections
func InitCassandraConnection(hosts []string, keyspace string, protocol int) {
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.ProtoVersion = protocol
	conn, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	cassandraConnection = conn
	log.Println("Connected to Cassandra Cluster")
}

// Following dhould not be called before InitCassandraConnection
// Also shouldn't be called directly
// Use the commands.go functions instead

func executeCassandraQuery(query string, values ...interface{}) {
	q := cassandraConnection.Query(query, values...)
	err := q.Exec()
	if err != nil {
		panic(err)
	}
}

func executeSelectCassandraQuery(query string, values ...interface{}) ([]map[string]interface{}, int) {
	var m []map[string]interface{}
	i := cassandraConnection.Query(query, values...).Iter()
	for {
		row := make(map[string]interface{})
		if !i.MapScan(row) {
			break
		}
		m = append(m, row)
	}
	return m, len(m)
}
