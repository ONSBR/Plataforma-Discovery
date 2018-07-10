package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "bankapp"
)

type Scan func(dest ...interface{}) error

func Query(binder func(Scan), query string, args ...interface{}) error {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)

	dataConn, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return err
	}
	defer dataConn.Close()
	dataConn.SetMaxOpenConns(10)
	dataConn.SetMaxIdleConns(10)
	result, err := dataConn.Query(query, args...)
	if err != nil {
		return err
	}
	for result.Next() {
		binder(result.Scan)
	}
	return nil
}
