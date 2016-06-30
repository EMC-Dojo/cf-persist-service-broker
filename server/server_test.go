package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/EMC-Dojo/cf-persist-service-broker/libstoragewrapper"
	"github.com/EMC-Dojo/cf-persist-service-broker/model"
	"github.com/EMC-Dojo/cf-persist-service-broker/utils"
	"github.com/emccode/libstorage/api/types"
)

var _ = Describe("Unit", func() {
	var serverURL, brokerUser, brokerPassword, instanceID, planIDString, libsHostServiceName, appGUID, bindingID, serviceBindingPath, driverType, storagePool string
	var libsClient types.APIClient

	type PlanID model.PlanID
	type ProvisionInstanceRequest model.ServiceInstance
	BeforeSuite(func() {

		//First Setup Variables we will use through the tests
		appGUID = "aaaa-bbbb-ccc-dddd"
		bindingID = "47E843FC-1A3A-4846-BC5D-E5F08BBD1CF1"
		storagePool = os.Getenv("STORAGE_POOL_NAME")
		Expect(storagePool).ToNot(BeEmpty())
		instanceID = os.Getenv("TEST_INSTANCE_ID") //3c653bce-8752-451b-96d9-a8a1a925b118
		Expect(instanceID).ToNot(BeEmpty())
		driverType = os.Getenv("LIBSTORAGE_DRIVER_TYPE")
		Expect(driverType).ToNot(BeEmpty())
		libsServerHost := os.Getenv("LIBSTORAGE_URI")
		Expect(libsServerHost).ToNot(BeEmpty())
		brokerUser = os.Getenv("BROKER_USERNAME")
		Expect(brokerUser).ToNot(BeEmpty())
		brokerPassword = os.Getenv("BROKER_PASSWORD")
		Expect(brokerPassword).ToNot(BeEmpty())
		port := os.Getenv("BROKER_PORT")
		Expect(port).ToNot(BeEmpty())
		insecureEnv := os.Getenv("INSECURE")
		Expect(insecureEnv).ToNot(BeEmpty())
		insecure, err := strconv.ParseBool(insecureEnv)
		Expect(err).ToNot(HaveOccurred())

		serverURL = "http://localhost:" + port
		serviceBindingPath = "/v2/service_instances/" + instanceID + "/service_bindings/" + bindingID
		//Second Start Local Server for tests!
		os.Chdir("..")
		s := Server{}

		devNull, err := os.Open(os.DevNull)
		if err != nil {
			log.Panic("Unable to open ", os.DevNull, err)
		}

		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = devNull
		gin.LoggerWithWriter(ioutil.Discard)

		s.Run(insecure, brokerUser, brokerPassword, port)
		//	time.Sleep(time.Millisecond * 500)

		libsClient = NewLibsClient()
		libsHostServiceName, err = libstoragewrapper.GetServiceNameByDriver(libsClient, driverType)
		Expect(err).ToNot(HaveOccurred())
		planIDString, err = utils.CreatePlanIDString(libsHostServiceName)
		Expect(err).ToNot(HaveOccurred())

	})

	AfterEach(func() {
		volumeID, err := libstoragewrapper.GetVolumeID(NewLibsClient(), libsHostServiceName, instanceID)
		if err == nil {
			libstoragewrapper.RemoveVolume(NewLibsClient(), libsHostServiceName, volumeID)
		}
	})

	Context("when fetching catalog", func() {
		It("returns catalog", func() {
			req, err := http.NewRequest("GET", serverURL+"/v2/catalog", nil)
			Expect(err).ToNot(HaveOccurred())
			req.SetBasicAuth(brokerUser, brokerPassword)
			resp, err := (&http.Client{}).Do(req)
			Expect(err).ToNot(HaveOccurred())
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			Expect(err).ToNot(HaveOccurred())

			var catalog model.Catalog
			err = json.Unmarshal(body, &catalog)
			Expect(err).ToNot(HaveOccurred())
			var serviceFromCatalog = catalog.Services[0]
			Expect(serviceFromCatalog).ToNot(BeNil())
			Expect(serviceFromCatalog.ID).To(Equal("c8ddac0a-36d3-41f7-bf72-990fe65b8d16"))
			Expect(serviceFromCatalog.Name).To(Equal("EMC-Persistence"))
			Expect(serviceFromCatalog.Description).ToNot(BeEmpty())
			Expect(serviceFromCatalog.Bindable).To(BeTrue())
			Expect(serviceFromCatalog.Requires[0]).To(Equal("volume_mount"))

			services, err := libstoragewrapper.GetServices(libsClient)
			Expect(err).To(BeNil())
			var plans = serviceFromCatalog.Plans
			Expect(plans).ToNot(BeNil())

			// iterates over all plans from response, ensures that the libstorage output is accounted for in service plans
			for _, plan := range plans {
				var planName model.PlanID
				err := json.Unmarshal([]byte(plan.ID), &planName)
				Expect(err).ToNot(HaveOccurred())
				var planMatchService bool
				for _, service := range services {
					if (planName.LibsServiceName == service.Name) && (planName.LibsHostName == os.Getenv("LIBSTORAGE_URI")) {
						planMatchService = true
						break
					}
				}
				Expect(planMatchService).To(BeTrue())
			}
			Expect(len(plans)).To(Equal(len(services)))

		})
	})

	It("provisions, binds, unbinds, and deprovisions", func() {
		// Provision
		testProvisionRequest := &ProvisionInstanceRequest{
			OrganizationGUID: "88011F5E-9FA0-484F-BFE5-9F1EED50B7D6",
			PlanID:           planIDString,
			ServiceID:        instanceID,
			SpaceGUID:        "A2331788-A736-479D-A9FB-114336F144C3",
			Parameters: map[string]string{
				"storage_pool_name": storagePool,
			},
			AcceptsIncomplete: true,
		}

		request, err := json.Marshal(testProvisionRequest)
		Expect(err).ToNot(HaveOccurred())

		path := "/v2/service_instances/" + instanceID
		req, err := http.NewRequest("PUT", serverURL+path, bytes.NewReader(request))
		Expect(err).ToNot(HaveOccurred())
		req.SetBasicAuth(brokerUser, brokerPassword)

		resp, err := (&http.Client{}).Do(req)
		Expect(err).ToNot(HaveOccurred())
		defer resp.Body.Close()
		Expect(resp.StatusCode).To(Equal(201))

		body, err := ioutil.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(strings.TrimSpace(string(body))).To(Equal("{}"))

		// Bind
		volumeID, err := libstoragewrapper.GetVolumeID(libsClient, libsHostServiceName, instanceID)
		Expect(err).ToNot(HaveOccurred())
		volumeName, err := utils.CreateNameForVolume(instanceID)

		expectedStructure := model.CreateServiceBindingResponse{
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

		provisionInstanceRequest := &model.ServiceBinding{
			ServiceID: libsHostServiceName,
			AppID:     appGUID,
			PlanID:    planIDString,
			BindResource: map[string]string{
				"app_guid": appGUID,
			},
		}

		provisionInstanceRequestBody, err := json.Marshal(provisionInstanceRequest)
		Expect(err).ToNot(HaveOccurred())

		req, err = http.NewRequest("PUT", serverURL+serviceBindingPath, bytes.NewReader(provisionInstanceRequestBody))
		Expect(err).ToNot(HaveOccurred())
		req.SetBasicAuth(brokerUser, brokerPassword)

		resp, err = (&http.Client{}).Do(req)
		Expect(err).ToNot(HaveOccurred())
		defer resp.Body.Close()
		Expect(resp.StatusCode).To(Equal(201))

		body, err = ioutil.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred())

		var binding model.CreateServiceBindingResponse
		err = json.Unmarshal(body, &binding)
		Expect(err).ToNot(HaveOccurred())
		Expect(binding).To(Equal(expectedStructure))

		// Remove binding
		u, err := url.Parse(serverURL + serviceBindingPath)
		q := u.Query()
		q.Set("service_id", instanceID)
		q.Set("plan_id", planIDString)
		u.RawQuery = q.Encode()
		req, err = http.NewRequest("DELETE", u.String(), nil)
		Expect(err).ToNot(HaveOccurred())
		req.SetBasicAuth(brokerUser, brokerPassword)

		resp, err = (&http.Client{}).Do(req)
		Expect(err).ToNot(HaveOccurred())
		defer resp.Body.Close()
		Expect(resp.StatusCode).To(Equal(200))

		body, err = ioutil.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(strings.TrimSpace(string(body))).To(Equal("{}"))

		// Removing instance
		u, err = url.Parse(serverURL + "/v2/service_instances/" + instanceID)
		q = u.Query()
		q.Set("service_id", libsHostServiceName)
		q.Set("plan_id", planIDString)
		u.RawQuery = q.Encode()
		req, err = http.NewRequest("DELETE", u.String(), nil)
		Expect(err).ToNot(HaveOccurred())

		req.SetBasicAuth(brokerUser, brokerPassword)
		resp, err = (&http.Client{}).Do(req)
		Expect(err).ToNot(HaveOccurred())
		defer resp.Body.Close()
		Expect(resp.StatusCode).To(Equal(200))

		body, err = ioutil.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(strings.TrimSpace(string(body))).To(Equal("{}"))
	})
})
