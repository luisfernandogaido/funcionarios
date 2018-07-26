package modelo

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"sync"
)

var (
	db   *sql.DB
	muRd sync.Mutex
	rd   redis.Conn
)

func Db(usuario, senha, servidor, banco string) error {
	dsn := fmt.Sprintf(
		"%v:%v@tcp(%v:3306)/%v?loc=America%%2FSao_Paulo&parseTime=true&multiStatements=true",
		usuario,
		senha,
		servidor,
		banco,
	)
	var err error
	db, err = sql.Open("mysql", dsn)
	return err
}

func Rd(addr string) error {
	var err error
	rd, err = redis.Dial("tcp", addr)
	return err
}
