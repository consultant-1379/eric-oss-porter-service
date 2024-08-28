package dbutils

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var F, e = os.OpenFile("/tmp/restsim.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
var PsqlInfo string
var Db *sql.DB
var Err error

// CreateTable creates a table in the specified database using the provided table name and columns.
func CreateTable(db *sql.DB, tableName string, columns []ColumnDefinition) error {
	// Construct the CREATE TABLE statement
	var columnsStrs []string
	for _, col := range columns {
		columnStr := fmt.Sprintf("%s %s", col.Name, col.Type)
		if col.PrimaryKey {
			columnStr += " PRIMARY KEY"
		}
		columnsStrs = append(columnsStrs, columnStr)
	}

	createStmt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, strings.Join(columnsStrs, ", "))

	// Execute the CREATE TABLE statement
	_, err := db.Exec(createStmt)
	if err != nil {
		return err
	}

	return nil
}

// ColumnDefinition represents a column definition in a database table.
type ColumnDefinition struct {
	Name       string
	Type       string
	PrimaryKey bool
}


// ExecuteSQLFromFile executes SQL statements from a file using the provided database connection.
func ExecuteSQLFromFile(db *sql.DB, filePath string) error {
	// Read the SQL file
	sqlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Split the file contents into individual SQL statements
	//statements := strings.Split(string(sqlFile), ";")

	stmt := string(sqlFile)

	_,err = db.Exec(stmt)
	if err != nil {
		return err
	}

	// Execute each SQL statement
	/*for _, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement != "" {
			_, err := db.Exec(statement)
			if err != nil {
				return err
			}
		}
	}*/

	return nil
}