package utils

import (
	"os"
	"path/filepath"
  "github.com/EMC-CMD/cf-persist-service-broker/model"
  "github.com/emccode/libstorage/api/types"
  "strings"
  "errors"
  "fmt"
  "github.com/emccode/libstorage/api/context"
)

func ProjectDirectory() string {
	return filepath.Join(os.Getenv("GOPATH"), "src/github.com/EMC-CMD/cf-persist-service-broker")
}

func GenerateVolumeName(instanceId string, serviceInstance model.ServiceInstance) (string, error) {
  // This will need to handle Isilon with a Conditional in the future, but for now, it just handles ScaleIO.

  return createNameForScaleIO(instanceId)
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

func GetVolumeID(libsClient types.Client, instanceId string, serviceInstance model.ServiceInstance) (string, error) {
  return getScaleIOVolumeID(libsClient, instanceId)
}

func getScaleIOVolumeID(libsClient types.Client, instanceID string) (string, error) {
  volumeName, err := createNameForScaleIO(instanceID)
  if err != nil {
    return "", fmt.Errorf("error creating name for volume using instance id %s", err)
  }

  volumeID, err := getScaleIOVolumeIDByName(libsClient, volumeName)
  if err != nil {
    return "", fmt.Errorf("error getting volume id from volume name %s", err)
  }

  return volumeID, nil
}

func getScaleIOVolumeIDByName(libsClient types.Client, volumeName string) (string, error) {
  ctx := context.Background()
  volumesOpts := types.VolumesOpts{
    Attachments: false,
  }

  volumes, err := libsClient.Storage().Volumes(ctx, &volumesOpts)
  if err != nil {
    return "", fmt.Errorf("error getting scaleio volumes %s", err)
  }

  for _, v := range volumes {
    if v.Name == volumeName {
      return v.ID, nil
    }
  }

  return "", fmt.Errorf("could not find volume id from volume name: %s", volumeName)
}

func createScaleIOVolumeParams(serviceInstance model.ServiceInstance) (az, volumeType string, iops, size int64, err error) {
  return "az", serviceInstance.Parameters["storage_pool_name"].(string), int64(100), int64(8), nil
}

func createNameForScaleIO(instanceId string) (string, error) {
  if len(instanceId) < 31 {
    return instanceId, nil
  }

  cooked_string := strings.Replace(instanceId, "-", "", -1)
  cooked_string = cooked_string[0:len(cooked_string)-1]

  if len(cooked_string) > 31 {
    return "", errors.New("Volume name cannot exceed 32 characters when all hyphens are removed.")
  }

  return cooked_string, nil
}