package main

import (
	"os"

	"github.com/EMC-Dojo/cf-persist-service-broker/server"
	log "github.com/Sirupsen/logrus"
)

func main() {
	s := server.Server{}
	username := os.Getenv("BROKER_USERNAME")
	password := os.Getenv("BROKER_PASSWORD")
	port := os.Getenv("PORT")

	insecure := (os.Getenv("INSECURE") == "true")

	log.Info("Starting EMC Persistance Service Broker at Port ", port)

	s.Run(insecure, username, password, port)
}
