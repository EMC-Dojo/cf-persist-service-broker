package main

import (
	"os"

	"flag"

	"github.com/EMC-Dojo/cf-persist-service-broker/server"
	log "github.com/Sirupsen/logrus"
	_ "github.com/emccode/libstorage/imports/executors"
	_ "github.com/emccode/libstorage/imports/remote"
	_ "github.com/emccode/libstorage/imports/routers"
)

func main() {
	s := server.Server{}
	configPath := flag.String("config", "", "Configuration override .yml file")
	flag.Parse()

	s.Init(*configPath)

	log.Info("Starting service at port ", os.Getenv("BROKER_PORT"))
	s.Run(os.Getenv("BROKER_PORT"))
}
