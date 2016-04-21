package server_test

import (
  "testing"
  "io/ioutil"

  log "github.com/Sirupsen/logrus"

  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
)

func TestServer(t *testing.T) {
  log.SetOutput(ioutil.Discard)

  RegisterFailHandler(Fail)
  RunSpecs(t, "Server Suite")
}
