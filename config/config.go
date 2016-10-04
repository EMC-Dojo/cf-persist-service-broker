package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/EMC-Dojo/cf-persist-service-broker/model"
)

func ReadConfig(configPath string) ([]model.Service, error) {
	configData, err := ioutil.ReadFile(configPath)
	services := []model.Service{}

	err = json.Unmarshal(configData, &services)
	if err != nil {
		return []model.Service{}, err
	}

	serviceID := services[0].ID
	if !validUUID(serviceID) {
		return services, fmt.Errorf("the service uuid given is not a valid uuid: %s", serviceID)
	}

	return services, nil
}

func validUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[8|9|aA|bB][a-f0-9]{3}-[a-f0-9]{12}$")
	return r.MatchString(uuid)
}
