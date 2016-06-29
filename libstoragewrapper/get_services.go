package libstoragewrapper

import (
	"fmt"
	"strings"

	"github.com/emccode/libstorage/api/context"
	"github.com/emccode/libstorage/api/types"
)

// GetServices : should return service names from libstorage catalog
func GetServices(libsClient types.APIClient) (map[string]*types.ServiceInfo, error) {
	services, err := libsClient.Services(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting service from LibStorage : (%s)", err)
	}
	return services, nil
}

// GetServiceNameByDriver : should return random service name by driverType from libstorage catalog
func GetServiceNameByDriver(libsClient types.APIClient, driverType string) (string, error) {
	services, err := GetServices(libsClient)
	if err != nil {
		return "", fmt.Errorf("error getting service for use with driver type %s : (%s)", driverType, err)
	}

	// run through and find a service by driver
	for servicename, service := range services {
		if strings.ToLower(service.Driver.Name) == strings.ToLower(driverType) {
			return servicename, nil
		}
	}

	return "", fmt.Errorf("Driver Type %s does not match service on libstorage", driverType)
}
