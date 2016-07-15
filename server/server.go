package server

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"

	"encoding/json"

	"github.com/EMC-Dojo/cf-persist-service-broker/libstoragewrapper"
	"github.com/EMC-Dojo/cf-persist-service-broker/model"
	"github.com/EMC-Dojo/cf-persist-service-broker/utils"
	"github.com/emccode/libstorage/api/client"
	"github.com/emccode/libstorage/api/types"
)

var serviceUUID, libStorageServiceName, libstorageHost, driverName string

// Server : Service Broker Server
type Server struct {
}

// NewLibsClient : creates a client used to communicate with libstorage.uri
func NewLibsClient() types.APIClient {
	insecure := (os.Getenv("INSECURE") == "true")
	libstorageHost = os.Getenv("LIBSTORAGE_URI")
	if libstorageHost == "" {
		log.Panic("A libstorage storage host must be specified with ENV[LIBSTORAGE_URI]")
	}
	return client.New(libstorageHost, &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	})
}

// Run the Service Broker
func (s Server) Run() {
	username := os.Getenv("BROKER_USERNAME")
	password := os.Getenv("BROKER_PASSWORD")
	port := os.Getenv("PORT")
	server := gin.Default()
	gin.SetMode("release")
	authorized := server.Group("/", gin.BasicAuth(gin.Accounts{
		username: password,
	}))

	libStorageServiceName = os.Getenv("LIB_STOR_SERVICE")
	driverName = os.Getenv("DIEGO_DRIVER_SPEC")
	serviceUUID = os.Getenv("EMC_SERVICE_UUID")

	fmt.Println("LibstorageServiceName: ", libStorageServiceName, "driverName", driverName, "serviceUUID", serviceUUID, "emcServiceName", os.Getenv("EMC_SERVICE_NAME"))
	r := regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[8|9|aA|bB][a-f0-9]{3}-[a-f0-9]{12}$")
	if !r.MatchString(serviceUUID) {
		log.Panic(fmt.Sprintf("The UUID given from ENV[EMC_SERVICE_UUID]= %s either was not set or is not a valid UUID", os.Getenv("EMC_SERVICE_UUID")))
	}

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
	// filter with supported drivers?
	services, err := libstoragewrapper.GetServices(NewLibsClient())
	if err != nil {
		log.Panicf("error retrieving services from libstorage host %s : (%s) ", os.Getenv("LIBSTORAGE_URI"), err)
	}
	var plans []model.Plan
	for _, service := range services {
		if service.Name == libStorageServiceName {
			planID, err := utils.CreatePlanIDString(service.Name)
			if err != nil {
				log.Panic(fmt.Sprintf("Error creating PlanID from name %s : (%s)", service.Name, err))
			}
			plans = append(plans, model.Plan{
				ID:          planID, // UUID made from JSON Marshalled Libstorage Host IP/service name
				Name:        service.Name,
				Description: service.Driver.Name,
			})
		}
	}
	if len(plans) < 1 {
		log.Panic(fmt.Sprintf("Service %s was not found on the LibStorage Server %s", libStorageServiceName, libstorageHost))
	}

	serviceName := os.Getenv("EMC_SERVICE_NAME")

	if serviceName == "" {
		serviceName = "EMC-Persistence"
	}

	catalogResponse := model.Catalog{
		Services: []model.Service{
			model.Service{
				ID:          serviceUUID,
				Name:        serviceName,
				Description: "Supports EMC ScaleIO & Isilon Storage Arrays for use with CloudFoundry",
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

	brokerName := os.Getenv("EMC_SERVICE_NAME")

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
				ContainerPath: fmt.Sprintf("/var/vcap/store/%s/%s", brokerName, volumeID),
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
