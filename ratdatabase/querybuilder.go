package ratdatabase

import (
	"bytes"
	"fmt"
	"strconv"
)

const (
	// Query Templates
	selectStatement = "SELECT %s FROM %s WHERE %s"
	insertStatement = "INSERT INTO %s (%s) VALUES (%s)"
	updateStatement = "UPDATE %s SET %s WHERE %s"
	deleteStatement = "DELETE FROM %s WHERE %s"
)

// Creates Field1, Field2....
func createCSVStringFromList(params []string) string {
	var commaSeparatedList bytes.Buffer

	for _, param := range params {
		commaSeparatedList.WriteString(param)
		commaSeparatedList.WriteString(",")
	}

	size := commaSeparatedList.Len() - 1 // Remove the last comma
	commaSeparatedList.Truncate(size)
	return commaSeparatedList.String()
}

// Creates Field=? and Field=? and ...
func createWhereStringFromList(params []string) string {
	var commaSeparatedList bytes.Buffer

	for _, val := range params {
		commaSeparatedList.WriteString(fmt.Sprintf(" %s=? and", val))
	}

	size := commaSeparatedList.Len() - 4 // Remove the last "and "
	commaSeparatedList.Truncate(size)
	return commaSeparatedList.String()
}

// Creates Field=?, Field=?, ...
func createUpdateStringFromList(params []string) string {
	var commaSeparatedList bytes.Buffer

	for _, val := range params {
		commaSeparatedList.WriteString(fmt.Sprintf(" %s=?,", val))
	}

	size := commaSeparatedList.Len() - 1 // Remove the last comma
	commaSeparatedList.Truncate(size)
	return commaSeparatedList.String()
}

// Creates a string of '?,?,?,?,...' to match the values column for Insert query
func createParameterString(size int) string {
	var commaSeparatedList bytes.Buffer
	for i := 0; i < size; i++ {
		commaSeparatedList.WriteString("?,")
	}
	// Every '?' is followed by a ',' so length is  2x the number of parameters.
	// -1 Because we dont want the last ','
	finalStringSize := (size * 2) - 1
	commaSeparatedList.Truncate(finalStringSize)
	return commaSeparatedList.String()
}

func stringArray(items ...string) []string {
	return items
}

func toString(number int) string {
	return strconv.Itoa(number)
}

func castInt(item interface{}) int {
	return int(item.(int64))
}

func castString(item interface{}) string {
	return item.(string)
}

func limitQuery(qry string, limit int) string {
	return fmt.Sprintf("%s LIMIT %d", qry, limit)
}

func createSelectQuery(columns []string, table string, conditions []string) string {
	cols := createCSVStringFromList(columns)
	conds := createWhereStringFromList(conditions)
	return fmt.Sprintf(selectStatement, cols, table, conds)
}

func createInsertStatement(table string, columns []string) string {
	cols := createCSVStringFromList(columns)
	vals := createParameterString(len(columns))
	return fmt.Sprintf(insertStatement, table, cols, vals)
}

func createUpdateStatement(table string, columns []string, cond []string) string {
	cols := createUpdateStringFromList(columns)
	conds := createWhereStringFromList(cond)
	return fmt.Sprintf(updateStatement, table, cols, conds)
}

func createDeleteStatement(table string, conditions []string) string {
	conds := createWhereStringFromList(conditions)
	return fmt.Sprintf(deleteStatement, table, conds)
}
