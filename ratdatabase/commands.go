package ratdatabase

const (
	// Query Templates
	selectStatement = "SELECT %s FROM %s WHERE %s"
	insertStatement = "INSERT INTO %s (%s) VALUES (%s)"
	updateStatement = "UPDATE %s SET %s WHERE %s"
	deleteStatement = "DELETE FROM %s where %s"

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
	userid = "userid"
	cash   = "usableCash"
)

var (
	// Database Helpers
	count = []string{"COUNT(*)"}
)

// Add

func UserExists(username string) bool {
	qry := createSelectQuery(count, users, stringArray(userid))
	rs := executeSelectCassandraQuery(qry, username)
	count := castInt(rs[0]["count"])
	return count == 1
}

func CreateUser(username string, cents int) {
	qry := createInsertStatement(users, stringArray(userid, cash))
	executeCassandraQuery(qry, username, toString(cents))
}

// Buy

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
