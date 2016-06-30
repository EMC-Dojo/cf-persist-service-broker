package server

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"encoding/json"

	"github.com/EMC-Dojo/cf-persist-service-broker/libstoragewrapper"
	"github.com/EMC-Dojo/cf-persist-service-broker/model"
	"github.com/EMC-Dojo/cf-persist-service-broker/utils"
	"github.com/emccode/libstorage/api/client"
	"github.com/emccode/libstorage/api/types"
)

var insecure bool

// Server : Service Broker Server
type Server struct {
}

// NewLibsClient : creates a client used to communicate with libstorage.uri
func NewLibsClient() types.APIClient {
	libstorageHost := os.Getenv("LIBSTORAGE_URI")
	if libstorageHost == "" {
		log.Panic("A libstorage storage host must be specified with ENV[LIBSTORAGE_URI]")
	}
	return client.New(libstorageHost, &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	})
}

// Run the Service Broker
func (s Server) Run(insecureEnv bool, username, password, port string) {
	insecure = insecureEnv
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
}

// CatalogHandler : Shows a catalog of available service plans
func CatalogHandler(c *gin.Context) {
	// filter with supported drivers?
	services, err := libstoragewrapper.GetServices(NewLibsClient())
	if err != nil {
		log.Panicf("error retrieving services from libstorage host %s : (%s) ", os.Getenv("LIBSTORAGE_URI"), err)
	}
	var plans []model.Plan
	for _, service := range services {
		planID, err := utils.CreatePlanIDString(service.Name)
		if err != nil {
			log.Panicf("Error creating PlanID from name %s : (%s)", service.Name, err)
		}
		plans = append(plans, model.Plan{
			ID:          planID, // UUID made from JSON Marshalled Libstorage Host IP/service name
			Name:        service.Name,
			Description: service.Driver.Name,
		})
	}

	catalogResponse := model.Catalog{
		Services: []model.Service{
			model.Service{
				ID:          "c8ddac0a-36d3-41f7-bf72-990fe65b8d16",
				Name:        "EMC-Persistence",
				Description: "EMC-Persistent Storage",
				Bindable:    true,
				Requires:    []string{"volume_mount"},
				Plans:       plans,
			},
		},
	}
	c.JSON(http.StatusOK, catalogResponse)
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

	volumeCreateRequest, err := utils.CreateVolumeRequest(instanceID, serviceInstance.Parameters["storage_pool_name"], int64(8))
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
	instanceID := c.Param("instanceID")
	volumeID, err := libstoragewrapper.GetVolumeID(NewLibsClient(), c.Query("service_id"), instanceID)
	if err != nil {
		log.Panicf("error service instance with ID %s does not exist. %s", instanceID, err)
	}

	err = libstoragewrapper.RemoveVolume(NewLibsClient(), c.Query("service_id"), volumeID)
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
		log.Panicf("Unable to read request body %s. Body", err)
	}

	var serviceBinding model.ServiceBinding
	err = json.Unmarshal(reqBody, &serviceBinding)
	if err != nil {
		log.Panicf("Unable to unmarshal service binding request %s. Request Body: %s", err, string(reqBody))
	}

	var planInfo = model.PlanID{}
	err = json.Unmarshal([]byte(serviceBinding.PlanID), &planInfo)
	if err != nil {
		log.Panic(fmt.Sprintf("Unable to unmarshal PlanID: %s", err))
	}
	serviceName := planInfo.LibsServiceName

	volumeID, err := libstoragewrapper.GetVolumeID(NewLibsClient(), serviceName, instanceID)
	if err != nil {
		log.Panicf("Unable to find volume ID by instance Id: %s", err)
	}

	volumeName, err := utils.CreateNameForVolume(instanceID)
	if err != nil {
		log.Panicf("Unable to encode instanceID to volume Name %s", err)
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
				ContainerPath: fmt.Sprintf("/var/vcap/store/scaleio/%s", volumeID),
				Mode:          "rw",
				Private: model.VolumeMountPrivateDetails{
					Driver:  "rexray",
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
