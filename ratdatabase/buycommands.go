package ratdatabase

import "github.com/gocql/gocql"

func InsertPendingBuyTransaction(username string, pendCash int, stockName string, stockPrice int) string {
	qry := createInsertStatement(pendingBuy, stringArray(pendingTID, userid, pendingCash, stock, stockValue))
	uuid := gocql.TimeUUID()
	executeCassandraQuery(qry, uuid, pendCash, stockName, stockPrice)
	return uuid.String()
}

func GetPendingBuyTransaction(username string, uuid string) (leftoverCash int, transExists bool) {
	qry := createSelectQuery(stringArray(pendingCash), pendingBuy, stringArray(pendingTID, userid))
	rs, count := executeSelectCassandraQuery(qry, uuid, username)

	if count == 0 {
		return // use 'zero' values as default (int=0, bool=false, string="")
	}
	leftoverCash = castInt(rs[0][pendingCash])
	transExists = true
	return
}

func DeletePendingBuyTransaction(username string, uuid string) {
	qry := createDeleteStatement(pendingBuy, stringArray(userid, pendingTID))
	executeCassandraQuery(qry, username, uuid)
}

func GetLastPendingBuyTransaction(username string) (uuid string, holdingCash int, stockName string, stockPrice int, exists bool) {
	qry := createSelectQuery(stringArray(pendingTID, pendingCash, stockname, stockValue), pendingBuy, stringArray(userid))
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
