package persistence

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDb(url string) (*sql.DB, error) {
	dbc, err := sql.Open("postgres", url)
	db = dbc
	return dbc, err
}
