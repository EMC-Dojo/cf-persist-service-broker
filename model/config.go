package model

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/akutz/gofig"
	"io"
	"os"
	"strconv"
)

var defaultConfigYaml = []byte(`
scaleio:
  insecure:             false
  useCerts:             true
  thinOrThick:          ThinProvisioned
  version:              2.0
`)

var Scaleio_required_configuration = []string{
  "scaleio.endpoint",
  "scaleio.insecure",
  "scaleio.useCerts",
  "scaleio.userName",
  "scaleio.password",
  "scaleio.systemID",
  "scaleio.systemName",
  "scaleio.protectionDomainID",
  "scaleio.protectionDomainName",
  "scaleio.storagePoolName",
  "scaleio.thinOrThick",
  "scaleio.version",
}

func GetConfig(in io.Reader) (gofig.Config, error) {
	c := gofig.New()
	return loadConfig(c, in)
}

func loadConfig(c gofig.Config, in io.Reader) (gofig.Config, error) {
	loadDefaults(c)

	err := overrideDefaultsFromIoReader(c, in)
	if err != nil {
		return c, err
	}

	err = overrideDefaultsFromEnv(c)
	if err != nil {
		return c, err
	}

	return c, validateConfiguration(c)
}

func loadDefaults(c gofig.Config) {
	c.ReadConfig(bytes.NewReader(defaultConfigYaml))
}

func overrideDefaultsFromIoReader(c gofig.Config, in io.Reader) error {
	err := c.ReadConfig(in)
	if err != nil {
		return err
	}

	return nil
}

func overrideDefaultsFromEnv(c gofig.Config) error {
	overrideIfProvided(c, "libstorage.storage.driver", "LIBSTORAGE_STORAGE_DRIVER")
	overrideIfProvided(c, "libstorage.host", "LIBSTORAGE_HOST")
	overrideIfProvided(c, "scaleio.endpoint", "SCALEIO_ENDPOINT")

	err := setBoolFromEnv(c, "scaleio.insecure", "SCALEIO_INSECURE")
	if err != nil {
		return err
	}
	err = setBoolFromEnv(c, "scaleio.usecerts", "SCALEIO_USE_CERTS")
	if err != nil {
		return err
	}

	overrideIfProvided(c, "scaleio.userName", "SCALEIO_USERNAME")
	overrideIfProvided(c, "scaleio.password", "SCALEIO_PASSWORD")
	overrideIfProvided(c, "scaleio.systemID", "SCALEIO_SYSTEM_ID")
	overrideIfProvided(c, "scaleio.systemName", "SCALEIO_SYSTEM_NAME")
	overrideIfProvided(c, "scaleio.protectionDomainID", "SCALEIO_PROTECTION_DOMAIN_ID")
	overrideIfProvided(c, "scaleio.protectionDomainName", "SCALEIO_PROTECTION_DOMAIN_NAME")
	overrideIfProvided(c, "scaleio.storagePoolName", "SCALEIO_STORAGE_POOL_NAME")
	overrideIfProvided(c, "scaleio.thinOrThick", "SCALEIO_THIN_OR_THICK")

	err = setFloatFromEnv(c, "scaleio.version", "SCALEIO_VERSION")
	if err != nil {
		return err
	}

	return nil
}

func validateConfiguration(c gofig.Config) error {
	if isEmpty(c, "libstorage.storage.driver") || isEmpty(c, "libstorage.host") {
		return errors.New("A libstorage storage driver and host must both be specified")
	}

	if read(c, "libstorage.storage.driver") != "scaleio" {
		return nil
	}

	missing_configuration := []string{}
	for _, v := range Scaleio_required_configuration {
		if isEmpty(c, v) {
			missing_configuration = append(missing_configuration, v)
		}
	}

	if len(missing_configuration) > 0 {
		return fmt.Errorf("Error validating configuration-missing necessary scaleio configuration %v", missing_configuration)
	}

	return nil
}

func overrideIfProvided(c gofig.Config, key string, variableName string) {
	value := readEnv(variableName)
	if value != "" {
		override(c, key, value)
	}
}

func override(c gofig.Config, key string, value interface{}) {
	c.Set(key, value)
}

func readEnv(variableName string) string {
	return os.Getenv(variableName)
}

func setBoolFromEnv(c gofig.Config, key string, variableName string) error {
	targetValue := readEnv(variableName)
	if targetValue != "" {
		result, err := strconv.ParseBool(targetValue)
		if err != nil {
			return err
		}
		override(c, key, result)
	}

	return nil
}

func setFloatFromEnv(c gofig.Config, key string, variableName string) error {
	targetValue := readEnv(variableName)
	if targetValue != "" {
		result, err := strconv.ParseFloat(targetValue, 64)
		if err != nil {
			return err
		}
		override(c, key, result)
	}

	return nil
}

func read(c gofig.Config, key string) interface{} {
	return c.Get(key)
}

func isEmpty(c gofig.Config, key string) bool {
	value := read(c, key)
	return value == "" || value == nil
}
