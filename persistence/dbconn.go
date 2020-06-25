package persistence

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

func dbConn() *sql.DB {
	c, err := sql.Open("postgres", os.Getenv("DB_CONN_URL"))
	if err != nil {
		panic(fmt.Sprint("db connection err", err))
	}
	return c
}
