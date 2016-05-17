package main

import (
	"os"

	"flag"
	"github.com/EMC-CMD/cf-persist-service-broker/server"
	log "github.com/Sirupsen/logrus"
)

func main() {
	s := server.Server{}
	configPath := flag.String("config", "", "Configuration override .yml file")
	flag.Parse()

	s.Init(*configPath)

	log.Info("Starting service at port ", os.Getenv("PORT"))
	s.Run(os.Getenv("PORT"))
}
