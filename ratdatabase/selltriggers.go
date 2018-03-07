package ratdatabase

import "github.com/gocql/gocql"

func SellTriggerExists(userID string, stockName string) (exists bool, amount int) {
	qry := createSelectQuery(stringArray(stockamount), selltrigger, stringArray(userid, stock))
	rs, count := executeSelectCassandraQuery(qry, userID, stockName)

	if count < 1 {
		return
	}

	exists = true
	amount = castInt(rs[0][stockamount])
	return
}

func InsertSetSellTrigger(userID string, stockName string, stockAmount int) (oldStockSellAmount int) {
	exists, oldStockSellAmount := SellTriggerExists(userID, stockName)

	if exists {
		qry := createUpdateStatement(selltrigger, stringArray(stockamount), stringArray(userid, stock))
		executeCassandraQuery(qry, stockAmount, userID, stockName)
	} else {
		qry := createInsertStatement(selltrigger, stringArray(pendingTTID, stockamount, stock, userid))
		uuid, _ := gocql.RandomUUID()
		executeCassandraQuery(qry, uuid, stockAmount, stockName, userID)
	}
	return
}

func UpdateSellTriggerPrice(userID string, stockName string, triggerValue int) (triggerAmountSet bool) {
	triggerAmountSet, _ = SellTriggerExists(userID, stockName)

	if triggerAmountSet {
		qry := createUpdateStatement(selltrigger, stringArray(triggervalue), stringArray(userid, stock))
		executeCassandraQuery(qry, triggerValue, userID, stockName)
	}

	return
}

func CancelSellTrigger(userID string, stockName string) (oldSellAmonut int) {
	_, oldSellAmonut = SellTriggerExists(userID, stockName)
	if oldSellAmonut == 0 {
		return
	}
	qry := createDeleteStatement(selltrigger, stringArray(userid, stock))
	executeCassandraQuery(qry, userID, stockName)
	return
}
