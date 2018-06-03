package postgres

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

type Executer interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	Select(dest interface{}, query string, args ...interface{}) error
}

const (
	host = "localhost"
)

func New() (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USERNAME"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DBNAME"))
	return sqlx.Open("postgres", psqlInfo)
}
