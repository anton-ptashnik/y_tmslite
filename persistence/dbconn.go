package persistence

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
)

var db *sql.DB

func InitDb() (*sql.DB, error) {
	dbc, err := sql.Open("postgres", os.Getenv("DB_CONN_URL"))
	db = dbc
	return dbc, err
}
