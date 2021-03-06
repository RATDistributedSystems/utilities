package ratdatabase

import (
	"github.com/gocql/gocql"
)

func InsertPendingBuyTransaction(username string, pendCash int, stockName string, stockPrice int) string {
	qry := createInsertStatement(pendingBuy, stringArray(pendingTID, userid, pendingcash, stock, stockValue))
	uuid := gocql.TimeUUID()
	executeCassandraQuery(qry, uuid, username, pendCash, stockName, stockPrice)
	return uuid.String()
}

func BuyTransactionAlive(username string, uuid string) bool {
	qry := createSelectQuery(countAll, pendingBuy, stringArray(pendingTID, userid))
	rs, length := executeSelectCassandraQuery(qry, uuid, username)
	numTransactions := 0
	if (length > 0){
		numTransactions = castInt64(rs[0][count])
	}
	return numTransactions == 1
}

func DeletePendingBuyTransaction(username string, uuid string) {
	qry := createDeleteStatement(pendingBuy, stringArray(userid, pendingTID))
	executeCassandraQuery(qry, username, uuid)
}

func GetLastPendingBuyTransaction(username string) (uuid string, holdingCash int, stockName string, stockPrice int, exists bool) {
	qry := createSelectQuery(stringArray(pendingTID, pendingcash, stock, stockValue), pendingBuy, stringArray(userid))
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
