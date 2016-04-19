package main

import (
	"os"

	"github.com/EMC-CMD/cf-persist-service-broker/server"
)

func main() {
	s := server.Server{}
	s.Run(os.Getenv("PORT"))
}
