// Package db provides connection to database server.
package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// New simply returns an instance of Database connection.
func New() (Db *sql.DB) {
	Db, err := sql.Open("mysql", "ba2f05c58da2e5:5b14e29b@tcp(us-cdbr-east-06.cleardb.net)/heroku_9a80e3a354317fc")
	if err != nil {
		log.Fatal("Failed to open")
	}

	err = Db.Ping()
	if err != nil {
		log.Fatal("Failed!")
	}
	return Db

}
