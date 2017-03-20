package models

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // used for sqlx
	"log"
	"os"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// InitDb initializes the database and returns the DB handle
func InitDb() *sqlx.DB {
	dbURI := os.Getenv("DATABASE_URI")
	if dbURI == "" {
		dbURI = "postgres://localhost/aws_analytics"
	}
	conn, err := sqlx.Open("postgres", dbURI)
	failOnError(err, "Unable to open connection to postgreSQL")
	return conn
}
