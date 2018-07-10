package main

import (
	"flag"
	"os"

	"github.com/ONSBR/Plataforma-Discovery/api"
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
	api.InitAPI()
}
