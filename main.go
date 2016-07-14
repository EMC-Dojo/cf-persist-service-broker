package main

import "github.com/EMC-Dojo/cf-persist-service-broker/server"

func main() {
	s := server.Server{}
	s.Run()
}
