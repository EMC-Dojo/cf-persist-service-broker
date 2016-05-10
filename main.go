package main

import (
  "os"

  log "github.com/Sirupsen/logrus"
  "github.com/EMC-CMD/cf-persist-service-broker/server"
  "flag"
)

func main() {
  s := server.Server{}
  configPath := flag.String("config", "", "Configuration override .yml file")
  flag.Parse()

  s.Init(*configPath)

  log.Info("Starting service at port ", os.Getenv("PORT"))
  s.Run(os.Getenv("PORT"))
}
