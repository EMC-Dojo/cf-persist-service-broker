package server

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"github.com/EMC-CMD/cf-persist-service-broker/model"
	"github.com/emccode/libstorage/api/types"
	"github.com/emccode/libstorage/client"
)

var scaleioClient types.Client

// The Service Broker Server
type Server struct {
}

func (s Server) SetClient(c types.Client) {
	scaleioClient = c
}

func (s Server) Init(configPath string) {
	if scaleioClient != nil {
		log.Info("client already set; skipping initialization")
		return
	}

	var configReader io.Reader
	var err error

	configReader = strings.NewReader("")
	if configPath != "" {
		configReader, err = os.Open(configPath)
		if err != nil {
			log.Panic("Unable to open ", configPath, err)
		}
	}

	config, err := model.GetConfig(configReader)
	if err != nil {
		log.Panic(err)
	}

	scaleioClient, err = client.New(config)
	if err != nil {
		log.Panic("Unable to create client", err)
	}
}

// Run the Service Broker
func (s Server) Run(port string) {
	server := gin.Default()
	gin.SetMode("release")
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
