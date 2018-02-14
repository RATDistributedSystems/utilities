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
	buyTrigger  = "buyTriggers"
	sellTrigger = "sellTriggers"

	// Audit Server Tables Names
	auditUserCommands       = "userCommands"
	auditQuoteRequests      = "quote_server"
	auditAccountTransaction = "account_transaction"
	auditSystemEvent        = "system_event"
	auditErrorEvent         = "error_event"
	auditDebugEvent         = "debug_event"

	// Common Field names
	userid      = "userid"
	balance     = "usableCash"
	pendingCash = "pendingCash"
	count       = "count"
	pendingTID  = "pid"
	stock       = "stock"
	stockValue  = "stockValue"
)

var (
	// Database Helpers
	countAll = []string{"COUNT(*)"}
)

// Add

func UserExists(username string) bool {
	qry := createSelectQuery(countAll, users, stringArray(userid))
	rs, _ := executeSelectCassandraQuery(qry, username)
	count := castInt(rs[0][count])
	return count == 1
}

func CreateUser(username string, cents int) {
	qry := createInsertStatement(users, stringArray(userid, balance))
	executeCassandraQuery(qry, username, toString(cents))
}

// Buy

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

func InsertPendingBuyTransaction(username string, pendCash int, stockName string, stockVal int) string {
	qry := createInsertStatement(pendingBuy, stringArray(pendingTID, userid, pendingCash, stock, stockValue))
	uuid := gocql.TimeUUID()
	executeCassandraQuery(qry, uuid, pendCash, stockName, stockVal)
	return uuid.String()
}

func GetPendingTransactionByUUID(username string, uuid string) (int, int) {
	qry := createSelectQuery(stringArray(pendingCash), pendingBuy, stringArray(pendingTID, userid))
	rs, count := executeSelectCassandraQuery(qry, uuid, username)
	if count == 0 {
		return 0, 0
	}
	leftoverCash := castInt(rs[0][pendingCash])
	return leftoverCash, count
}

func DeletePendingTransaction(username string, uuid string) {
	qry := createDeleteStatement(pendingBuy, stringArray(userid, pendingTID))
	executeCassandraQuery(qry, username, uuid)
}

// Buy Cancel

// Buy Commit

// Sell

// Sell Cancel

// Sell Commit

// Buy Set Amount

// Buy Trigger

// Buy Trigger Cancel

// Sell set amount

// Sell Trigger

// Sell Trigger cancel

// Display summary

// Quote Command

// Sell Command
