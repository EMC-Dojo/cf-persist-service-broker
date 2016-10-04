package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/EMC-Dojo/cf-persist-service-broker/model"
)

var DefaultConfig = []model.Service{
	model.Service{
		ID:          "92e98925-d046-4c72-9598-ba352449a5c7",
		Name:        "Persistent-Storage",
		Description: "Supports EMC ScaleIO & Isilon Storage Arrays for use with CloudFoundry",
		Bindable:    true,
		Requires:    []string{"volume_mount"},
		Plans: []model.Plan{
			model.Plan{
				Name:        "isilonservice",
				Description: "An isilon service",
			},
			model.Plan{
				Name:        "scaleioservice",
				Description: "A scaleio service",
			},
		},
	},
}

func ReadConfig(configPath string) ([]model.Service, error) {
	if configPath == "" {
		return DefaultConfig, nil
	}

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
