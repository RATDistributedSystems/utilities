package ratdatabase

import (
	"fmt"

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
	userid      = "userid"
	balance     = "usablecash"
	pendingcash = "pendingcash"
	count       = "count"
	pendingTID  = "pid"
	userstockid = "usid"
	stock       = "stock"
	stockValue  = "stockValue"
	stockname   = "stockname"
	stockamount = "stockamount"
)

var (
	// Database Helpers
	countAll = []string{"COUNT(*)"}
)

func UserExists(username string) bool {
	qry := createSelectQuery(countAll, users, stringArray(userid))
	rs, _ := executeSelectCassandraQuery(qry, username)
	count := castInt64(rs[0][count])
	return count == 1
}

func CreateUser(username string, cents int) {
	qry := createInsertStatement(users, stringArray(userid, balance))
	executeCassandraQuery(qry, username, toString(cents))
}

func GetUserBalance(username string) int {
	qry := createSelectQuery(stringArray(balance), users, stringArray(userid))
	rs, count := executeSelectCassandraQuery(qry, username)
	if count != 1 {
		panic(fmt.Sprintf("No user '%s' found\n", username))
	}

	balance := castInt(rs[0][balance])
	return balance
}

func UpdateUserBalance(username string, newBalance int) {
	qry := createUpdateStatement(users, stringArray(balance), stringArray(userid))
	executeCassandraQuery(qry, newBalance, username)
}

func GetStockAmountOwned(username string, stockName string) (uuid string, stockAmount int, exists bool) {
	qry := createSelectQuery(stringArray(stockname, stockamount, userstockid), userstocks, stringArray(userid))
	rs, _ := executeSelectCassandraQuery(qry, username)

	for _, r := range rs {
		if r[stockname] == stockName {
			stockAmount = castInt(r[stockamount])
			uuid = castString(r[userstockid])
			exists = true
			return
		}
	}

	uuid = ""
	stockAmount = 0
	exists = false
	return
}

func AddStockToPortfolio(username string, stockName string, stockAmount int) {
	qry := createInsertStatement(userstocks, stringArray(userstockid, userid, stockamount, stockname))
	uuid, _ := gocql.RandomUUID()
	executeCassandraQuery(qry, uuid, username, stockAmount, stockName)
}

func UpdateUserStockByUUID(uuid string, stockName string, stockAmount int) {
	qry := createUpdateStatement(userstocks, stringArray(stockamount), stringArray(userstockid))
	executeCassandraQuery(qry, stockAmount, uuid)
}

// Buy Set Amount

/// Sell

// Sell Commit/ Buy Trigger

// Buy Trigger Cancel

// Sell set amount

// Sell Trigger

// Sell Trigger cancel

// Display summary

// Quote Command

// Sell Command
