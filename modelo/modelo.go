package modelo

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"strings"
)

var db *sql.DB

func Db(dsn string) error {
	var err error
	dsn = strings.Replace(dsn, `"`, "", 2)
	db, err = sql.Open("mysql", dsn)
	return err
}
