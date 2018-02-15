package ratdatabase

import (
	"github.com/gocql/gocql"
)

func InsertPendingSellTransaction(username string, stockName string, pendCash string, stockVal int) string {
	cols := stringArray(pendingTID, userid, pendingCash, stock, stockValue)
	qry := createInsertStatement(pendingSell, cols)
	uuid := gocql.TimeUUID()
	executeCassandraQuery(qry, uuid, username, pendCash, stockName, stockVal)
	return uuid.String()
}

func GetPendingSellTransaction(username string, uuid string) (saleProfits int, stockVal int, exists bool) {
	qry := createSelectQuery(stringArray(pendingCash, stockValue), pendingSell, stringArray(pendingTID, userid))
	rs, count := executeSelectCassandraQuery(qry, uuid, username)

	if count == 0 { // No pending transactions
		return
	}
	saleProfits = castInt(rs[0][pendingCash])
	stockVal = castInt(rs[0][stockValue])
	exists = true
	return
}

func DeletePendingSellTransaction(username string, uuid string) {
	qry := createDeleteStatement(pendingSell, stringArray(userid, pendingTID))
	executeCassandraQuery(qry, username, uuid)
}

func GetLastPendingSellTransaction(username string) (uuid string, holdingCash int, stockName string, stockPrice int, exists bool) {
	qry := createSelectQuery(stringArray(pendingTID, pendingCash, stockname, stockValue), pendingSell, stringArray(userid))
	qry = limitQuery(qry, 1)
	rs, count := executeSelectCassandraQuery(qry, username)

	if count == 0 { // no record to return
		return
	}
	uuid = castString(rs[0][pendingTID])
	holdingCash = castInt(rs[0][pendingCash])
	stockName = castString(rs[0][stockname])
	stockPrice = castInt(rs[0][stockValue])
	exists = true
	return
}
