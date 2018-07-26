package main

import (
	"log"
	"os"

	"github.com/luisfernandogaido/funcionarios/modelo"
	"github.com/luisfernandogaido/funcionarios/server"
)

func main() {
	if err := modelo.Db(os.Getenv("MYSQL")); err != nil {
		log.Fatal(err)
	}
	log.Fatal(server.Start(":4003"))
}
