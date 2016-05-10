package utils

import (
  "os"
  "path/filepath"
)

func ProjectDirectory() string {
  return filepath.Join(os.Getenv("GOPATH"), "src/github.com/EMC-CMD/cf-persist-service-broker")
}
