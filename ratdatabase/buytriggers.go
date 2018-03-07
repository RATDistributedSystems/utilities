package ratdatabase

import (
	"github.com/gocql/gocql"
)

func BuyTriggerExists(userID string, stockName string) (exists bool, amount int) {
	qry := createSelectQuery(stringArray(pendingcash), buytrigger, stringArray(userid, stock))
	rs, count := executeSelectCassandraQuery(qry, userID, stockName)

	if count < 1 {
		return
	}

	exists = true
	amount = castInt(rs[0][pendingcash])
	return
}

func InsertSetBuyTrigger(userID string, stockName string, cashAmount int) (oldAmount int) {
	exists, oldAmount := BuyTriggerExists(userID, stockName)

	if exists {
		qry := createUpdateStatement(buytrigger, stringArray(pendingcash), stringArray(userid, stock))
		executeCassandraQuery(qry, cashAmount, userID, stockName)
	} else {
		qry := createInsertStatement(buytrigger, stringArray(pendingTTID, pendingcash, stock, userid))
		uuid, _ := gocql.RandomUUID()
		executeCassandraQuery(qry, uuid, cashAmount, stockName, userID)
	}
	return
}

func UpdateBuyTriggerPrice(userID string, stockName string, triggerPrice int) (triggerAmountSet bool) {
	triggerAmountSet, _ = BuyTriggerExists(userID, stockName)

	if triggerAmountSet {
		qry := createUpdateStatement(buytrigger, stringArray(triggervalue), stringArray(userid, stock))
		executeCassandraQuery(qry, triggerPrice, userID, stockName)
	}

	return
}

func CancelBuyTrigger(userID string, stockName string) (buyAmount int) {
	_, buyAmount = BuyTriggerExists(userID, stockName)
	if buyAmount == 0 {
		return
	}
	qry := createDeleteStatement(buytrigger, stringArray(userid, stock))
	executeCassandraQuery(qry, userID, stockName)
	return
}
