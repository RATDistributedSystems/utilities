package ratdatabase

import (
	"github.com/gocql/gocql"
)

func InsertPendingSellTransaction(username string, stockName string, pendCash int, stockVal int) string {
	cols := stringArray(pendingTID, userid, pendingcash, stock, stockValue)
	qry := createInsertStatement(pendingSell, cols)
	uuid := gocql.TimeUUID()
	executeCassandraQuery(qry, uuid, username, pendCash, stockName, stockVal)
	return uuid.String()
}

func SellTransactionAlive(username string, uuid string) bool {
	qry := createSelectQuery(countAll, pendingSell, stringArray(pendingTID, userid))
	rs, _ := executeSelectCassandraQuery(qry, uuid, username)
	numTransactions := castInt64(rs[0][count])
	return numTransactions == 1
}

func DeletePendingSellTransaction(username string, uuid string) {
	qry := createDeleteStatement(pendingSell, stringArray(userid, pendingTID))
	executeCassandraQuery(qry, username, uuid)
}

func GetLastPendingSellTransaction(username string) (uuid string, holdingCash int, stockName string, stockPrice int, exists bool) {
	qry := createSelectQuery(stringArray(pendingTID, pendingcash, stock, stockValue), pendingSell, stringArray(userid))
	qry = limitQuery(qry, 1)
	rs, count := executeSelectCassandraQuery(qry, username)

	if count == 0 { // no record to return
		return
	}
	uuid = castUUID(rs[0][pendingTID])
	holdingCash = castInt(rs[0][pendingcash])
	stockName = castString(rs[0][stock])
	stockPrice = castInt(rs[0][stockValue])
	exists = true
	return
}
