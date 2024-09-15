package config

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Connect() {
	database, err := sql.Open("mysql", "root:test@/bookmanagmentsystemdb")
	if err != nil {
		log.Fatal(err)
	}
	db = database
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
}
func GetDb() *sql.DB {
	return db
}
