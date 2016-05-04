package main

import (
  "os"

  log "github.com/Sirupsen/logrus"
  "github.com/EMC-CMD/cf-persist-service-broker/server"
)

func main() {
  if len(os.Args) < 2 {
    log.Panic("configuration for client is required")
  }
  s := server.Server{}
  s.Init(os.Args[1])

  log.Info("Starting service at port ", os.Getenv("PORT"))
  s.Run(os.Getenv("PORT"))
}
