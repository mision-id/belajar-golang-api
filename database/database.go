package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(ConnectionString string) (*sql.DB, error) {
	//Open Database Connection
	db, err := sql.Open("postgres", ConnectionString)
	if err != nil {
		return nil, err
	}

	// Test Datbase Connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	//set connection pool setting (optional tapi recomended)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Database connected Successfully")
	return db, nil
}
