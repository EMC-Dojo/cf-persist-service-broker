package utils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/EMC-CMD/cf-persist-service-broker/model"
	"github.com/emccode/libstorage/api/types"
)

func ProjectDirectory() string {
	return filepath.Join(os.Getenv("GOPATH"), "src/github.com/EMC-CMD/cf-persist-service-broker")
}

func GenerateVolumeName(instanceId string, serviceInstance model.ServiceInstance) (string, error) {
	// This will need to handle Isilon with a Conditional in the future, but for now, it just handles ScaleIO.

	return CreateNameForScaleIO(instanceId)
}

func CreateVolumeOpts(serviceInstance model.ServiceInstance) (types.VolumeCreateOpts, error) {
	// This will need to handle Isilon with a Conditional in the future, but for now, it just handles ScaleIO.
	var az, volumeType string
	var iops, size int64
	var err error

	az, volumeType, iops, size, err = createScaleIOVolumeParams(serviceInstance)

	volumeCreateOpts := types.VolumeCreateOpts{
		AvailabilityZone: &az,
		IOPS:             &iops,
		Size:             &size,
		Type:             &volumeType,
	}

	return volumeCreateOpts, err
}

func createScaleIOVolumeParams(serviceInstance model.ServiceInstance) (az, volumeType string, iops, size int64, err error) {
	return "az", serviceInstance.Parameters["storage_pool_name"], int64(100), int64(8), nil
}

func CreateNameForScaleIO(instanceId string) (string, error) {
	if len(instanceId) < 31 {
		return instanceId, nil
	}

	cooked_string := strings.Replace(instanceId, "-", "", -1)
	cooked_string = cooked_string[0 : len(cooked_string)-1]

	if len(cooked_string) > 31 {
		return "", errors.New("Volume name cannot exceed 32 characters when all hyphens are removed.")
	}

	return cooked_string, nil
}
