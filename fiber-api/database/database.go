package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

var (
	Connection *sql.DB
)

func init() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 55000, "postgres", "postgrespw", "postgres")
	var err error
	Connection, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}
