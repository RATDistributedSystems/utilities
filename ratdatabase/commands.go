package ratdatabase

import (
	//"fmt"

	"github.com/gocql/gocql"
)

const (
	// Transaction Server Table Names
	triggers    = "triggers"
	users       = "users"
	userstocks  = "userstocks"
	pendingBuy  = "buypendingtransactions"
	pendingSell = "sellpendingtransactions"
	buytrigger  = "buytriggers"
	selltrigger = "selltriggers"

	// Audit Server Tables Names
	auditusercommands       = "userCommands"
	auditquoterequests      = "quote_server"
	auditaccounttransaction = "account_transaction"
	auditsystemevent        = "system_event"
	auditrrrorevent         = "error_event"
	auditdebugevent         = "debug_event"

	// Common Field names
	userid       = "userid"
	balance      = "usablecash"
	pendingcash  = "pendingcash"
	count        = "count"
	pendingTID   = "pid"
	pendingTTID  = "tid"
	userstockid  = "usid"
	stock        = "stock"
	stockValue   = "stockvalue"
	stockamount  = "stockamount"
	triggervalue = "triggervalue"
)

var (
	// Database Helpers
	countAll = []string{"COUNT(*)"}
)

func UserExists(username string) bool {
	qry := createSelectQuery(stringArray(userid), users, stringArray(userid))
	_, count := executeSelectCassandraQuery(qry, username)
	return count == 1
}

func CreateUser(username string, cents int) {
	qry := createInsertStatement(users, stringArray(userid, balance))
	executeCassandraQuery(qry, username, toString(cents))
}

func GetUserBalance(username string) int {
	qry := createSelectQuery(stringArray(balance), users, stringArray(userid))
	rs, _ := executeSelectCassandraQuery(qry, username)
	/*
	if count != 1 {
		fmt.Println(rs[0])
		fmt.Sprintf("No user '%s' found\n", username)
	}
	*/
	balance := castInt(rs[0][balance])
	return balance
}

func UpdateUserBalance(username string, newBalance int) {
	qry := createUpdateStatement(users, stringArray(balance), stringArray(userid))
	executeCassandraQuery(qry, newBalance, username)
}

func GetStockAmountOwned(username string, stockName string) (uuid string, stockAmount int, exists bool) {
	qry := createSelectQuery(stringArray(stock, stockamount, userstockid), userstocks, stringArray(userid,stock))
	rs, _ := executeSelectCassandraQuery(qry, stringArray(username,stockName))

	for _, r := range rs {
		if r[stock] == stockName {
			stockAmount = castInt(r[stockamount])
			uuid = castUUID(r[userstockid])
			exists = true
			return
		}
	}

	return
}

func AddStockToPortfolio(username string, stockName string, stockAmount int) {
	qry := createInsertStatement(userstocks, stringArray(userstockid, userid, stockamount, stock))
	uuid, _ := gocql.RandomUUID()
	executeCassandraQuery(qry, uuid, username, stockAmount, stockName)
}

func UpdateUserStockByUUID(uuid string, stockName string, stockAmount int) {
	qry := createUpdateStatement(userstocks, stringArray(stockamount), stringArray(userstockid))
	executeCassandraQuery(qry, stockAmount, uuid)
}

func UpdateUserStockByUserAndStock(username string, stockName string, stockAmount int) {
	qry := createUpdateStatement(userstocks, stringArray(stockamount), stringArray(userid, stock))
	executeCassandraQuery(qry, stockAmount, username, stockName)
}

// Display summary

func GetStockAndAmountOwned(username string) ([]map[string]interface{}) {
	qry := createSelectQuery(stringArray(stock, stockamount), userstocks, stringArray(userid))
	rs, _ := executeSelectCassandraQuery(qry, username)

	return rs
}

func GetSellTriggers(username string) ([]map[string]interface{}) {
	qry := createSelectQuery(stringArray(stock,stockamount), selltrigger, stringArray(userid))
	rs, _ := executeSelectCassandraQuery(qry, username)
	
	return rs
}

func GetBuyTriggers(username string) ([]map[string]interface{}){
	qry := createSelectQuery(stringArray(stock,stockamount), buytrigger, stringArray(userid))
	rs, _ := executeSelectCassandraQuery(qry, username)

	return rs
}
