package database

import (
	"database/sql"
	"os"

	//Mysql Driver
	_ "github.com/go-sql-driver/mysql"
)

var (
	dbUser = os.Getenv("DB_USERNAME")
	dbPass = os.Getenv("DB_PASSWORDE")
	dbHost = os.Getenv("DB_HOST")
	dbName = os.Getenv("DB_NAME")
)

//DbConn returns connection to database
func DbConn() (db *sql.DB) {
	db, err := sql.Open(dbHost, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}
