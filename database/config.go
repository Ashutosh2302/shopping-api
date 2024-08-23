package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	connString := os.Getenv("DATABASE_URL")
	DB, err := sql.Open("postgres", connString)
	if err != nil {
		panic("could not connect to database")
	}
	if err = DB.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("The database is connected")
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	return DB
}
