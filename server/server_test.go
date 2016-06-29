package server

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/EMC-Dojo/cf-persist-service-broker/libstoragewrapper"
	"github.com/EMC-Dojo/cf-persist-service-broker/model"
	"github.com/EMC-Dojo/cf-persist-service-broker/utils"
	"github.com/emccode/libstorage/api/client"
	"github.com/emccode/libstorage/api/types"
)

var _ = Describe("Unit", func() {
	var serverURL string
	var brokerUser string
	var brokerPassword string
	var libsClient types.APIClient
	var instanceID string
	var planIDString string
	var libsHostServiceName string
	var appGUID string
	var bindingID string
	var serviceBindingPath string
	type PlanID model.PlanID
	type ProvisionInstanceRequest model.ServiceInstance

	BeforeEach(func() {
		var err error
		appGUID = "aaaa-bbbb-ccc-dddd"
		instanceID = os.Getenv("TEST_INSTANCE_ID") //3c653bce-8752-451b-96d9-a8a1a925b118
		bindingID = "47E843FC-1A3A-4846-BC5D-E5F08BBD1CF1"
		libsHostServiceName = "josh"
		serviceBindingPath = "/v2/service_instances/" + instanceID + "/service_bindings/" + bindingID

		planIDString, err = utils.CreatePlanIDString(libsHostServiceName)
		Expect(err).To(BeNil())
		libsServerHost := os.Getenv("LIBSTORAGE_URI")
		Expect(libsServerHost).ToNot(BeEmpty())

		port := os.Getenv("BROKER_PORT")
		Expect(port).ToNot(BeEmpty())
		serverURL = "http://localhost:" + port

		brokerUser = os.Getenv("BROKER_USERNAME")
		Expect(serverURL).ToNot(BeEmpty())
		brokerPassword = os.Getenv("BROKER_PASSWORD")
		Expect(serverURL).ToNot(BeEmpty())

		libsClient = client.New(libsServerHost, &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		})
		AllowInsecureConnections()
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
				"storage_pool_name": "default",
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
