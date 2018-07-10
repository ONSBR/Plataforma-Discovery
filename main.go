package main

import (
	"flag"
	"os"

	"github.com/ONSBR/Plataforma-Discovery/db"
	"github.com/labstack/gommon/log"
)

var local bool

func init() {
	flag.BoolVar(&local, "local", false, "to run service with local configuration")
}

func main() {
	flag.Parse()
	log.SetLevel(log.DEBUG)
	if local {
		os.Setenv("PORT", "8090")
	}
	//api.InitAPI()
	type conta struct {
		Id    string
		Saldo int
	}
	contas := make([]conta, 0)
	db.Query(func(scan func(dest ...interface{}) error) {
		var c conta
		scan(&c.Id, &c.Saldo)
		contas = append(contas, c)
	}, "select id, saldo from conta")
}
