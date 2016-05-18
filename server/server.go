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

	"encoding/json"

	"github.com/EMC-CMD/cf-persist-service-broker/libstoragewrapper"
	"github.com/EMC-CMD/cf-persist-service-broker/model"
	"github.com/EMC-CMD/cf-persist-service-broker/utils"
	"github.com/emccode/libstorage/api/context"
	"github.com/emccode/libstorage/api/types"
	"github.com/emccode/libstorage/client"
)

var libsClient types.Client

const (
	scaleio_service = "c8ddac0a-36d3-41f7-bf72-990fe65b8d16"
)

const (
	small_plan = "92798c7d-e7b0-49d6-8872-4aeafbb193ef"
)

// The Service Broker Server
type Server struct {
}

func (s Server) Init(configPath string) {
	if libsClient != nil {
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

	libsClient, err = client.New(config)
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
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Panic("Unable to read request body %s", err)
	}

	var serviceInstance model.ServiceInstance
	err = json.Unmarshal(reqBody, &serviceInstance)
	if err != nil {
		log.Panicf("Unable to unmarshal the request body: %s. Request body %s", err, string(reqBody))
	}

	instanceId := c.Param("instanceId")
	volumeName, err := utils.GenerateVolumeName(instanceId, serviceInstance)
	if err != nil {
		log.Panic("Unable to generate volume name: %s.", err)
	}

	volumeOpts, err := utils.CreateVolumeOpts(serviceInstance)
	if err != nil {
		log.Panic("Unable to create volume opts: %s.", err)
	}

	ctx := context.Background()
	_, err = libstoragewrapper.CreateVolume(libsClient, ctx, volumeName, volumeOpts)
	if err != nil {
		log.Panic("Unable to create volume using libstorage: %s.", err)
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func DeprovisionHandler(c *gin.Context) {
	serviceInstance := model.ServiceInstance{
		ServiceId: c.Param("service_id"),
		PlanId:    c.Param("plan_id"),
	}

	instanceId := c.Param("instanceId")
	volumeID, err := libstoragewrapper.GetVolumeID(libsClient, instanceId, serviceInstance)
	if err != nil {
		log.Panic("Unable to find volume ID by instance Id")
	}

	ctx := context.Background()
	err = libstoragewrapper.RemoveVolume(libsClient, ctx, volumeID)
	if err != nil {
		log.Panic("error removing volume using libstorage")
	}

	c.JSON(http.StatusOK, gin.H{})
}

func BindingHandler(c *gin.Context) {
	body, _ := ioutil.ReadFile("fixtures/create_binding_response.json")
	c.String(http.StatusCreated, string(body))
}

func UnbindingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
