package database

import (
	"database/sql"
	"fmt"
	"os"

	//Mysql Driver
	_ "github.com/go-sql-driver/mysql"
)

var (
	dbUser   = os.Getenv("DB_USERNAME")
	dbPass   = os.Getenv("DB_PASSWORD")
	dbName   = os.Getenv("DB_NAME")
	dbURL    = "tcp(mysql-0.mysql:3306)"
	dbDriver = "mysql"
)

//DbConn returns connectiosn to database
func DbConn() (db *sql.DB) {
	db, err := sql.Open(dbDriver, fmt.Sprintf("%s:%s@%s/%s", dbUser, dbPass, dbURL, dbName))
	if err != nil {
		panic(err.Error())
	}
	return db
}
