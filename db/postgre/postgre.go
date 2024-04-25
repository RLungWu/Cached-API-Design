package postgre

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func CreatePGClient() *sql.DB {
	connStr := os.Getenv("POSTGRESQL_URL")
	if connStr == "" {
		log.Fatal("POSTGRESQL_URL is not set")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	err = db.Ping()
	if err != nil{
		log.Fatalf("Error pinging database: %q", err)
	}

	fmt.Println("Successfully connected to PostgreSQL")
	return db
}