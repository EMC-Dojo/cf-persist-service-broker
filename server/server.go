package server

import (
	"os"
	"github.com/gin-gonic/gin"
  "path/filepath"
  "net/http"
  "io/ioutil"

  "github.com/emccode/libstorage/client"
  "github.com/akutz/gofig"
  log "github.com/Sirupsen/logrus"
)

// The Service Broker Server
type Server struct {
}

func clientConfig() client.Client {
  config := gofig.New()
  file, err := os.Open("../config/config.yml")
  if err != nil {
    log.Panic("Unable to open config.yml", err)
  }
  config.ReadConfig(file)

  scaleioClient, err := client.New(config)
  if err != nil {
    log.Panic("Unable to create client", err)
  }

  return scaleioClient
}

// Run the Service Broker
func (s Server) Run(port string) {
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()
	authorized := server.Group("/", gin.BasicAuth(gin.Accounts{
		os.Getenv("BROKER_USERNAME"): os.Getenv("BROKER_PASSWORD"),
	}))

	authorized.GET("/v2/catalog", CatalogHandler)
	authorized.PUT("/v2/service_instances/:instanceId", ProvisioningHandler)
	authorized.PUT("/v2/service_instances/:instanceId/service_bindings/:bindingId", BindingHandler)
	authorized.DELETE("/v2/service_instances/:instanceId/service_bindings/:bindingId", UnbindingHandler)
	authorized.DELETE("/v2/service_instances/:instanceId", DeprovisionHandler)

	server.Run(":" + port)
}

func CatalogHandler(c *gin.Context) {
	c.Status(http.StatusOK)
	p, _ := filepath.Abs("templates/catalog.json")
	c.File(p)
}

func ProvisioningHandler(c *gin.Context) {
  c.JSON(http.StatusCreated, gin.H{})
}

func DeprovisionHandler(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{})
}

func BindingHandler(c *gin.Context) {
  body, _ := ioutil.ReadFile("fixtures/create_binding_response.json")
  c.String(http.StatusCreated, string(body))
}

func UnbindingHandler(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{})
}
