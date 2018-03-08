package ratdatabase

import (
	"log"
	"strconv"

	"github.com/gocql/gocql"
)

var CassandraConnection *gocql.Session

// InitCassandraConnection is called by a server to initialize the connections
func InitCassandraConnection(host string, keyspace string, protocol string) {
	cluster := gocql.NewCluster(host)
	cluster.Keyspace = keyspace
	cluster.ConnectTimeout = time.Second * 10
	proto, err := strconv.Atoi(protocol)
	if err != nil {
		panic(err)
	}
	cluster.ProtoVersion = proto
	conn, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	CassandraConnection = conn
	log.Println("Connected to Cassandra Cluster")
}

// Following dhould not be called before InitCassandraConnection
// Also shouldn't be called directly
// Use the commands.go functions instead

func executeCassandraQuery(query string, values ...interface{}) {
	q := CassandraConnection.Query(query, values...)
	err := q.Exec()
	if err != nil {
		panic(err)
	}
}

func executeSelectCassandraQuery(query string, values ...interface{}) ([]map[string]interface{}, int) {
	var m []map[string]interface{}
	i := CassandraConnection.Query(query, values...).Iter()
	for {
		row := make(map[string]interface{})
		if !i.MapScan(row) {
			break
		}
		m = append(m, row)
	}

	return m, len(m)
}
