package main

import (
	"log"
	"os"

	"github.com/luisfernandogaido/funcionarios/modelo"
	"github.com/luisfernandogaido/funcionarios/server"
)

func main() {
	if err := modelo.Db(
		os.Getenv("MYSQL_USUARIO"),
		os.Getenv("MYSQL_SENHA"),
		os.Getenv("MYSQL_SERVIDOR"),
		os.Getenv("MYSQL_BANCO"),
	); err != nil {
		log.Fatal(err)
	}
	if err := modelo.Rd(os.Getenv("REDIS")); err != nil {
		log.Fatal(err)
	}
	if err := modelo.SetMd(
		os.Getenv("MONGODB_ADDRS"),
		os.Getenv("MONGODB_DATABASE"),
		os.Getenv("MONGODB_USERNAME"),
		os.Getenv("MONGODB_PASSWORD"),
	); err != nil {
		log.Fatal(err)
	}
	server.Start(":" + os.Getenv("PORTA"))
}
