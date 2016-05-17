package server

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestServer(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Suite")
}

var _ = BeforeSuite(func() {
	go startServer()
	time.Sleep(time.Millisecond * 500)
})

func startServer() {
	os.Chdir("..")
	os.Setenv("PORT", "9900")
	s := Server{}
	devNull, err := os.Open(os.DevNull)
	if err != nil {
		log.Panic("Unable to open ", os.DevNull, err)
	}

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = devNull
	gin.LoggerWithWriter(ioutil.Discard)

	s.Init("")
	s.Run("9900")
}
