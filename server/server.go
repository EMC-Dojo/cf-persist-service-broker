package server

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"encoding/json"

	"github.com/EMC-Dojo/cf-persist-service-broker/config"
	"github.com/EMC-Dojo/cf-persist-service-broker/libstoragewrapper"
	"github.com/EMC-Dojo/cf-persist-service-broker/model"
	"github.com/EMC-Dojo/cf-persist-service-broker/utils"
	"github.com/emccode/libstorage/api/client"
	"github.com/emccode/libstorage/api/types"
)

var serviceUUID, libstorageHost, driverName string

// Server : Service Broker Server
type Server struct {
}

// NewLibsClient : creates a client used to communicate with libstorage.uri
func NewLibsClient() types.APIClient {
	insecure := (os.Getenv("INSECURE") == "true")
	libstorageHost = os.Getenv("LIBSTORAGE_URI")
	return client.New(libstorageHost, &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	})
}

// Run the Service Broker
func (s Server) Run() {
	expectingENVs := []string{"BROKER_USERNAME", "BROKER_PASSWORD", "PORT", "DIEGO_DRIVER_SPEC", "LIBSTORAGE_URI", "INSECURE"}
	err := VerifyEnvironmentVariable(expectingENVs)
	if err != nil {
		log.Panicf("error: %s. expecting envs: %s", err, expectingENVs)
	}

	username := os.Getenv("BROKER_USERNAME")
	password := os.Getenv("BROKER_PASSWORD")
	port := os.Getenv("PORT")
	driverName = os.Getenv("DIEGO_DRIVER_SPEC")

	server := gin.Default()
	gin.SetMode("release")
	authorized := server.Group("/", gin.BasicAuth(gin.Accounts{
		username: password,
	}))

	authorized.GET("/v2/catalog", CatalogHandler)
	authorized.PUT("/v2/service_instances/:instanceID", ProvisioningHandler)
	authorized.PUT("/v2/service_instances/:instanceID/service_bindings/:bindingId", BindingHandler)
	authorized.DELETE("/v2/service_instances/:instanceID/service_bindings/:bindingId", UnbindingHandler)
	authorized.DELETE("/v2/service_instances/:instanceID", DeprovisionHandler)

	server.Run(":" + port)
	log.Info("Starting EMC Persistance Service Broker at Port ", port)
}

// CatalogHandler : Shows a catalog of available service plans
func CatalogHandler(c *gin.Context) {
	configPath := os.Getenv("BROKER_CONFIG_PATH")
	catalogServices, err := config.ReadConfig(configPath)
	if err != nil {
		log.Panicf("error reading config file %s", err)
	}

	// filter with supported drivers?
	libstorageServices, err := libstoragewrapper.GetServices(NewLibsClient())
	if err != nil {
		log.Panicf("error retrieving services from libstorage host %s : (%s) ", os.Getenv("LIBSTORAGE_URI"), err)
	}

	plansExists := DoPlansExistInLibstorage(catalogServices[0].Plans, libstorageServices)
	if !plansExists {
		log.Panic("plan(s) do not exist in libstorage services.")
	}

	catalogServices[0].Plans, err = AddCatalogPlanIDs(catalogServices[0].Plans, libstorageHost)
	if err != nil {
		log.Panic("could not modify plans' ids")
	}

	c.JSON(http.StatusOK, model.Catalog{Services: catalogServices})
}

// ProvisioningHandler : Provisions service instances (e.g. creates volumes; used by cloud controller)
func ProvisioningHandler(c *gin.Context) {
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Panic(fmt.Sprintf("Unable to read request body %s", err))
	}

	var serviceInstance model.ServiceInstance
	err = json.Unmarshal(reqBody, &serviceInstance)
	if err != nil {
		log.Panicf("Unable to unmarshal the request body: %s. Request body %s", err, string(reqBody))
	}

	instanceID := c.Param("instanceID")
	var planInfo = model.PlanID{}
	err = json.Unmarshal([]byte(serviceInstance.PlanID), &planInfo)
	if err != nil {
		log.Panic(fmt.Sprintf("Unable to unmarshal PlanID: %s", err))
	}
	serviceName := planInfo.LibsServiceName

	volumeSize := int64(8)
	if serviceInstance.Parameters.SizeInGB != "" {
		volumeSize, err = strconv.ParseInt(serviceInstance.Parameters.SizeInGB, 10, 64)
		if err != nil {
			log.Panicf("Invalid SizeInGB %s", serviceInstance.Parameters.SizeInGB)
		}
	}

	volumeCreateRequest, err := utils.CreateVolumeRequest(instanceID, serviceInstance.Parameters.StoragePoolName, int64(volumeSize))
	if err != nil {
		log.Panic(fmt.Sprintf("Unable to create volume request: %s.", err))
	}

	_, err = libstoragewrapper.CreateVolume(NewLibsClient(), serviceName, volumeCreateRequest)
	if err != nil {
		log.Panic(fmt.Sprintf("Unable to create volume %+v using libstorage: %s.", volumeCreateRequest, err))
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// DeprovisionHandler : Deprovisions service instances (e.g. destroys volumes; used by cloud controller)
func DeprovisionHandler(c *gin.Context) {
	var planInfo = model.PlanID{}
	err := json.Unmarshal([]byte(c.Query("plan_id")), &planInfo)
	if err != nil {
		log.Panic(fmt.Sprintf("Unable to unmarshal PlanID: %s", err))
	}
	serviceName := planInfo.LibsServiceName

	instanceID := c.Param("instanceID")
	volumeID, err := libstoragewrapper.GetVolumeID(NewLibsClient(), serviceName, instanceID)
	if err != nil {
		log.Panic(fmt.Sprintf("error service instance with ID %s does not exist. %s", instanceID, err))
	}

	err = libstoragewrapper.RemoveVolume(NewLibsClient(), serviceName, volumeID)
	if err != nil {
		log.Panic("error removing volume using libstorage", err)
	}

	c.JSON(http.StatusOK, gin.H{})
}

// BindingHandler : binds volumes to applications (used by cloud controller)
func BindingHandler(c *gin.Context) {
	instanceID := c.Param("instanceID")

	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Panic(fmt.Sprintf("Unable to read request body %s. Body", err))
	}

	var serviceBinding model.ServiceBinding
	err = json.Unmarshal(reqBody, &serviceBinding)
	if err != nil {
		log.Panic(fmt.Sprintf("Unable to unmarshal service binding request %s. Request Body: %s", err, string(reqBody)))
	}

	var planInfo = model.PlanID{}
	err = json.Unmarshal([]byte(serviceBinding.PlanID), &planInfo)
	if err != nil {
		log.Panic(fmt.Sprintf("Unable to unmarshal PlanID: %s", err))
	}
	serviceName := planInfo.LibsServiceName

	volumeID, err := libstoragewrapper.GetVolumeID(NewLibsClient(), serviceName, instanceID)
	if err != nil {
		log.Panic(fmt.Sprintf("Unable to find volume ID by instance Id: %s", err))
	}

	volumeName, err := utils.CreateNameForVolume(instanceID)
	if err != nil {
		log.Panic(fmt.Sprintf("Unable to encode instanceID to volume Name %s", err))
	}

	serviceBindingResp := model.CreateServiceBindingResponse{
		Credentials: model.CreateServiceBindingCredentials{
			Database: "dummy",
			Host:     "dummy",
			Password: "dummy",
			Port:     3306,
			URI:      "dummy",
			Username: "dummy",
		},
		VolumeMounts: []model.VolumeMount{
			model.VolumeMount{
				//should we be using volumeID?
				ContainerPath: fmt.Sprintf("/var/vcap/store/%s", volumeID),
				Mode:          "rw",
				Private: model.VolumeMountPrivateDetails{
					Driver:  driverName,
					GroupId: volumeName,
					Config:  "{\"broker\":\"specific_values\"}",
				},
			},
		},
	}
	c.JSON(http.StatusCreated, serviceBindingResp)
}

// UnbindingHandler : unbinds volumes from applications (used by cloud controller)
func UnbindingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
