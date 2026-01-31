package database

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func InitDB(connectionString string) (*sql.DB, error) {
	// open connection to database
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, err
	}

	// verify connection is working
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings (optional but recommended)
	db.SetMaxOpenConns(25) // maximum number of open connections
	db.SetMaxIdleConns(5)  // maximum number of idle connections

	log.Println("Database connected successfully")
	return db, nil
}
