package main

import (
	"os"
	"github.com/EMC-CMD/cf-persist-service-broker/server"

	log "github.com/Sirupsen/logrus"
)

func main() {
	s := server.Server{}
	log.Info("Starting service at port ", os.Getenv("PORT"))
	s.Run(os.Getenv("PORT"))
}
