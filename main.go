package main

import (
	"log"
	"os"

	"github.com/luisfernandogaido/funcionarios/modelo"
	"github.com/luisfernandogaido/funcionarios/server"
)

func main() {
	err := modelo.Db(
		os.Getenv("MYSQL_USUARIO"),
		os.Getenv("MYSQL_SENHA"),
		os.Getenv("MYSQL_SERVIDOR"),
		os.Getenv("MYSQL_BANCO"),
	)
	if err != nil {
		log.Fatal(err)
	}
	if err = modelo.Rd(os.Getenv("REDIS")); err != nil {
		log.Fatal(err)
	}
	server.Start(":" + os.Getenv("PORTA"))
}
