package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	connStr := "host=postgres port=5432 user=telemetry_user password=telemetry_pass dbname=telemetry_db sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	return DB.Ping()
}