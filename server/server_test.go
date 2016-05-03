package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/EMC-CMD/cf-persist-service-broker/model"
	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"time"
)

var s *exec.Cmd

func startServer() *exec.Cmd {
	err := os.Chdir("..")
	Expect(err).ToNot(HaveOccurred())

	server := exec.Command("go", "run", "main.go")
	server.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	server.Stdout = ioutil.Discard
	server.Stderr = ioutil.Discard

	err = server.Start()
	Expect(err).ToNot(HaveOccurred())

	err = os.Chdir("server")
	Expect(err).ToNot(HaveOccurred())

	return server
}

var _ = BeforeSuite(func() {
	os.Setenv("PORT", "9900")
	err := godotenv.Load("../test.env")
	Expect(err).ToNot(HaveOccurred())
	s = startServer()

	time.Sleep(5 * time.Second)
})

var _ = AfterSuite(func() {
	s.Process.Signal(syscall.SIGHUP)
	pgid, err := syscall.Getpgid(s.Process.Pid)
	if err == nil {
		syscall.Kill(-pgid, 15) // note the minus sign
	}
	s.Wait()
	Expect(s.ProcessState.String()).To(Or(Equal("signal: terminated"), Equal("exit status 0")))
})

var _ = Describe("Integration", func() {

	var serverURL string
	var brokerUser string
	var brokerPassword string

	BeforeEach(func() {
		serverURL = "http://localhost:" + os.Getenv("PORT")
		Expect(serverURL).ToNot(BeEmpty())
		brokerUser = os.Getenv("BROKER_USERNAME")
		Expect(serverURL).ToNot(BeEmpty())
		brokerPassword = os.Getenv("BROKER_PASSWORD")
		Expect(serverURL).ToNot(BeEmpty())
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

			var expectedCatalog model.Catalog
			expectedBody, err := ioutil.ReadFile("../templates/catalog.json")
			Expect(err).ToNot(HaveOccurred())

			err = json.Unmarshal(expectedBody, &expectedCatalog)
			Expect(err).ToNot(HaveOccurred())
			Expect(catalog).To(Equal(expectedCatalog))
		})
	})

	Context("when provisioning instances", func() {
		Context("when request is valid", func() {
			It("returns 201 created", func() {
				provisionInstanceRequestBody, err := os.Open("../fixtures/provision_instance_request.json")
				Expect(err).ToNot(HaveOccurred())

				path := "/v2/service_instances/29C39AEB-9A09-49D3-A432-AE995C75FFF8"
				req, err := http.NewRequest("PUT", serverURL+path, provisionInstanceRequestBody)
				Expect(err).ToNot(HaveOccurred())
				req.SetBasicAuth(brokerUser, brokerPassword)

				resp, err := (&http.Client{}).Do(req)
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(201))

				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())
				Expect(strings.TrimSpace(string(body))).To(Equal("{}"))
			})
		})
	})

	Context("when creating bindings", func() {
		Context("when request is valid", func() {
			It("returns the binding authorization parameters with status 201", func() {
				provisionInstanceRequestBody, err := os.Open("../fixtures/create_binding_request.json")
				Expect(err).ToNot(HaveOccurred())

				path := "/v2/service_instances/CCDB8015-92BE-42FB-B4C3-00CEAB503144/service_bindings/47E843FC-1A3A-4846-BC5D-E5F08BBD1CF1"
				req, err := http.NewRequest("PUT", serverURL+path, provisionInstanceRequestBody)
				Expect(err).ToNot(HaveOccurred())
				req.SetBasicAuth(brokerUser, brokerPassword)

				resp, err := (&http.Client{}).Do(req)
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(201))

				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())

				var binding model.ServiceBinding
				err = json.Unmarshal(body, &binding)
				Expect(err).ToNot(HaveOccurred())

				expectedBody, err := ioutil.ReadFile("../fixtures/create_binding_response.json")
				Expect(err).ToNot(HaveOccurred())

				var expectedBinding model.ServiceBinding
				err = json.Unmarshal(expectedBody, &expectedBinding)
				Expect(err).ToNot(HaveOccurred())
				Expect(binding).To(Equal(expectedBinding))
			})
		})
	})

	Context("when removing bindings", func() {
		Context("when request is valid", func() {
			It("returns 200 ok", func() {
				path := "/v2/service_instances/CCCBCA4D-FFA9-4FA5-BF71-9584F7DB149F/service_bindings/07D68661-C6CC-4B3E-9991-F62C1AA3AAC6"
				u, err := url.Parse(serverURL + path)
				q := u.Query()
				q.Set("service_id", "1C12FB88-2F67-4708-8AB7-4215B8E27C3E")
				q.Set("plan_id", "205F2EF0-2B83-492F-9840-F585D3D8D6B8")
				u.RawQuery = q.Encode()
				req, err := http.NewRequest("DELETE", u.String(), nil)
				Expect(err).ToNot(HaveOccurred())
				req.SetBasicAuth(brokerUser, brokerPassword)

				resp, err := (&http.Client{}).Do(req)
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(200))

				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())
				Expect(strings.TrimSpace(string(body))).To(Equal("{}"))
			})
		})
	})

	Context("when removing instances", func() {
		Context("when request is valid", func() {
			It("returns 200 ok", func() {
				u, err := url.Parse(serverURL + "/v2/service_instances/BFF0C0CB-B811-4E4E-8930-F04377BF43C9")
				q := u.Query()
				q.Set("service_id", "1C12FB88-2F67-4708-8AB7-4215B8E27C3E")
				q.Set("plan_id", "205F2EF0-2B83-492F-9840-F585D3D8D6B8")
				u.RawQuery = q.Encode()
				req, err := http.NewRequest("DELETE", u.String(), nil)
				Expect(err).ToNot(HaveOccurred())

				req.SetBasicAuth(brokerUser, brokerPassword)
				resp, err := (&http.Client{}).Do(req)
				Expect(err).ToNot(HaveOccurred())
				defer resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(200))

				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())
				Expect(strings.TrimSpace(string(body))).To(Equal("{}"))
			})
		})
	})
})
