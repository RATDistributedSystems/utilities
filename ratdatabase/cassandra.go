package ratdatabase

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gocql/gocql"
)

var CassandraConnection *gocql.Session

// InitCassandraConnection is called by a server to initialize the connections
func InitCassandraConnection(host string, keyspace string, protocol string) {
	hostNoSpace := strings.TrimSpace(host)
	hosts := strings.Split(hostNoSpace, ",")
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.ConnectTimeout = time.Second * 1
	proto, err := strconv.Atoi(protocol)
	if err != nil {
		panic(err)
	}
	cluster.ProtoVersion = proto
	cluster.NumConns = 2 //number of connections per host
	cluster.RetryPolicy = 0 //retry policy for queries
	cluster.MaxPreparedStmts = 1000 //max cache size for prepared statemetns
	cluster.MaxRoutingKeyInfo = 1000 //max cache size for query info about statements for each session
	cluster.PageSize = 5000 //page size to use for created sessions

	for i := 0; i < 10; i++ {
		conn, err := cluster.CreateSession()
		if err != nil {
			log.Print("Audit DB not up yet. Waiting 10 more seconds...")
			time.Sleep(time.Second * 10)
			continue
		}
		CassandraConnection = conn
		log.Println("Connected to Cassandra Cluster")
		return
	}

	log.Fatalf("Couldn't connect to Cassandra Cluster %s", hostNoSpace)
}

// Following dhould not be called before InitCassandraConnection
// Also shouldn't be called directly
// Use the commands.go functions instead

func executeCassandraQuery(query string, values ...interface{}) {
	//error with new update here
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
