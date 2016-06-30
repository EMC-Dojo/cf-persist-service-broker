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

	username := os.Getenv("BROKER_USERNAME")
	Expect(username).ToNot(BeEmpty())
	password := os.Getenv("BROKER_PASSWORD")
	Expect(password).ToNot(BeEmpty())
	port := os.Getenv("BROKER_PORT")
	Expect(port).ToNot(BeEmpty())
	insecure := os.Getenv("INSECURE")
	Expect(insecure).ToNot(BeEmpty())

	log.Info("Starting EMC Persistance Service Broker at Port %s", port)

	s.Run(insecure, username, password, port)
}
