package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func ConnectDB() {
	var (
		dbDriver = "mysql"
		dbSource = "root:password@tcp(mysql:3306)/ordermanagmentsystem_db"
	)
	database, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal(err)
	}

	db = database

	log.Println("Connected to db ")
}
func GetDB() *sql.DB {
	return db
}
