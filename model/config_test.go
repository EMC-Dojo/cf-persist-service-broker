package model_test

import (
	. "github.com/EMC-CMD/cf-persist-service-broker/model"

	"github.com/EMC-CMD/cf-persist-service-broker/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
	"strings"
)

var _ = Describe("Config", func() {
	var configurationYaml = filepath.Join(utils.ProjectDirectory(), "config/config_test.yml")
	BeforeEach(func() {
		varNames := []string{
			"LIBSTORAGE_HOST",
			"LIBSTORAGE_STORAGE_DRIVER",
			"SCALEIO_ENDPOINT",
			"SCALEIO_INSECURE",
			"SCALEIO_USE_CERTS",
			"SCALEIO_USERNAME",
			"SCALEIO_PASSWORD",
			"SCALEIO_SYSTEM_ID",
			"SCALEIO_SYSTEM_NAME",
			"SCALEIO_PROTECTION_DOMAIN_ID",
			"SCALEIO_PROTECTION_DOMAIN_NAME",
			"SCALEIO_STORAGE_POOL_NAME",
			"SCALEIO_THIN_OR_THICK",
			"SCALEIO_VERSION",
		}
		for _, v := range varNames {
			err := os.Unsetenv(v)
			Expect(err).ToNot(HaveOccurred())
		}
	})
	It("Should raise an error when trying to use it without a driver or host", func() {
		_, err := GetConfig(strings.NewReader(""))
		Expect(err).To(MatchError("A libstorage storage driver and host must both be specified"))
	})

	It("Should raise an error when trying to use scaleio with libstorage without proper configuration", func() {
		os.Setenv("LIBSTORAGE_STORAGE_DRIVER", "scaleio")
		os.Setenv("LIBSTORAGE_HOST", "libstorage.com")
		_, err := GetConfig(strings.NewReader(""))
		Expect(err).To(MatchError("Error validating configuration-missing necessary scaleio configuration [scaleio.endpoint scaleio.userName scaleio.password scaleio.systemID scaleio.systemName scaleio.protectionDomainID scaleio.protectionDomainName scaleio.storagePoolName]"))
	})

	It("Has the right defaults", func() {
		os.Setenv("LIBSTORAGE_HOST", "https://libstorage.com/api")
		os.Setenv("LIBSTORAGE_STORAGE_DRIVER", "scaleio")
		os.Setenv("SCALEIO_ENDPOINT", "abc-d")
		os.Setenv("SCALEIO_SYSTEM_ID", "abc-d")
		os.Setenv("SCALEIO_USERNAME", "abc-d")
		os.Setenv("SCALEIO_PASSWORD", "abc-d")
		os.Setenv("SCALEIO_SYSTEM_NAME", "cluster1")
		os.Setenv("SCALEIO_PROTECTION_DOMAIN_ID", "abc-d")
		os.Setenv("SCALEIO_PROTECTION_DOMAIN_NAME", "abc-d")
		os.Setenv("SCALEIO_STORAGE_POOL_NAME", "abc-d")
		c, err := GetConfig(strings.NewReader(""))

		Expect(err).ToNot(HaveOccurred())
		Expect(c.Get("scaleio.insecure")).To(Equal(false))
		Expect(c.Get("scaleio.useCerts")).To(Equal(true))
		Expect(c.Get("scaleio.thinOrThick")).To(Equal("ThinProvisioned"))
		Expect(c.Get("scaleio.version")).To(Equal(2.0))

	})

	It("Overrides default values from a config file", func() {
		fileReader, err := os.Open(configurationYaml)
		Expect(err).ToNot(HaveOccurred())
		c, err := GetConfig(fileReader)

		Expect(c.Get("libstorage.host")).To(Equal("tcp://file_fake_host:9000"))
		Expect(c.Get("libstorage.storage.driver")).To(Equal("scaleio"))
		Expect(c.Get("scaleio.endpoint")).To(Equal("https://file_fake_endpoint/api"))
		Expect(c.Get("scaleio.insecure")).To(Equal(true))
		Expect(c.Get("scaleio.useCerts")).To(Equal(false))
		Expect(c.Get("scaleio.userName")).To(Equal("file_fake_user"))
		Expect(c.Get("scaleio.password")).To(Equal("file_fake_password"))
		Expect(c.Get("scaleio.systemID")).To(Equal("file_fake_sys_id"))
		Expect(c.Get("scaleio.systemName")).To(Equal("file_fake_sys_name"))
		Expect(c.Get("scaleio.protectionDomainID")).To(Equal("file_fake_protection_domain_id"))
		Expect(c.Get("scaleio.protectionDomainName")).To(Equal("file_fake_protection_domain_name"))
		Expect(c.Get("scaleio.storagePoolName")).To(Equal("file_fake_storage_pool_name"))
		Expect(c.Get("scaleio.thinOrThick")).To(Equal("ThinProvisioned"))
		Expect(c.Get("scaleio.version")).To(Equal(2.0))
	})

	It("Overrides file-configured values using the environment", func() {
		os.Setenv("LIBSTORAGE_HOST", "https://env_fake_host/api")
		os.Setenv("LIBSTORAGE_STORAGE_DRIVER", "env_fake_driver")
		os.Setenv("SCALEIO_ENDPOINT", "https://env_fake_endpoint/api")
		os.Setenv("SCALEIO_INSECURE", "false")
		os.Setenv("SCALEIO_USE_CERTS", "true")
		os.Setenv("SCALEIO_USERNAME", "env_fake_username")
		os.Setenv("SCALEIO_PASSWORD", "env_fake_password")
		os.Setenv("SCALEIO_SYSTEM_ID", "env_fake_system_id")
		os.Setenv("SCALEIO_SYSTEM_NAME", "env_fake_system_name")
		os.Setenv("SCALEIO_PROTECTION_DOMAIN_ID", "env_fake_protection_domain_id")
		os.Setenv("SCALEIO_PROTECTION_DOMAIN_NAME", "env_fake_protection_domain_name")
		os.Setenv("SCALEIO_STORAGE_POOL_NAME", "env_fake_storage_pool_name")
		os.Setenv("SCALEIO_THIN_OR_THICK", "ThickProvisioned")
		os.Setenv("SCALEIO_VERSION", "999.0")

		fileReader, err := os.Open(configurationYaml)
		Expect(err).ToNot(HaveOccurred())
		c, err := GetConfig(fileReader)
		Expect(c.Get("libstorage.host")).To(Equal("https://env_fake_host/api"))
		Expect(c.Get("libstorage.storage.driver")).To(Equal("env_fake_driver"))
		Expect(c.Get("scaleio.endpoint")).To(Equal("https://env_fake_endpoint/api"))
		Expect(c.Get("scaleio.insecure")).To(Equal(false))
		Expect(c.Get("scaleio.useCerts")).To(Equal(true))
		Expect(c.Get("scaleio.userName")).To(Equal("env_fake_username"))
		Expect(c.Get("scaleio.password")).To(Equal("env_fake_password"))
		Expect(c.Get("scaleio.systemID")).To(Equal("env_fake_system_id"))
		Expect(c.Get("scaleio.systemName")).To(Equal("env_fake_system_name"))
		Expect(c.Get("scaleio.protectionDomainID")).To(Equal("env_fake_protection_domain_id"))
		Expect(c.Get("scaleio.protectionDomainName")).To(Equal("env_fake_protection_domain_name"))
		Expect(c.Get("scaleio.storagePoolName")).To(Equal("env_fake_storage_pool_name"))
		Expect(c.Get("scaleio.thinOrThick")).To(Equal("ThickProvisioned"))
		Expect(c.Get("scaleio.version")).To(Equal(999.0))
	})
})
