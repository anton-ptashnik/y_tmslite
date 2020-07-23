package persistence

import (
	"github.com/joho/godotenv"
	"os"
)

func init()  {
	godotenv.Load("../.env")
	_, err := InitDb(os.Getenv("TESTDB_CONN_URL"))
	panicOnErr(err)
}