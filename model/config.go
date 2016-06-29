package model

import (
	"log"
	"os"
)

var config map[string]string

func init() {
	config = map[string]string{}
	libstorageHost := os.Getenv("LIBSTORAGE_URI")
	if libstorageHost == "" {
		log.Panic("A libstorage storage host must be specified")
	}
	config["libstorage.uri"] = libstorageHost
}

// GetConfig : Provides access to configuration for cf-persist-service-broker
func GetConfig() map[string]string {
	return config
}
