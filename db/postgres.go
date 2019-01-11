package db

import (
	"database/sql"
	"fmt"

	"github.com/ONSBR/Plataforma-Deployer/env"
	"github.com/ONSBR/Plataforma-Discovery/helpers"

	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
)

type Scan func(dest ...interface{}) error

func Query(systemID string, binder func(Scan), query string, args ...interface{}) error {
	dbName, err := helpers.GetDBName(systemID)
	if err != nil {
		log.Error(err)
		return err
	}
	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		env.Get("POSTGRES_HOST", "localhost"), env.Get("POSTGRES_PORT", "5432"), env.Get("POSTGRES_USER", DB_USER), env.Get("POSTGRES_PASSWORD", DB_PASSWORD), dbName)

	dataConn, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Error(err)
		return err
	}
	defer dataConn.Close()
	dataConn.SetMaxOpenConns(10)
	dataConn.SetMaxIdleConns(10)
	result, err := dataConn.Query(query, args...)
	if err != nil {
		log.Error(err)
		return err
	}
	for result.Next() {
		binder(result.Scan)
	}
	return nil
}
