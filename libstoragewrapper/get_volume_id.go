package libstoragewrapper

import (
	"fmt"

	"github.com/emccode/libstorage/api/context"

	"github.com/EMC-CMD/cf-persist-service-broker/utils"
	"github.com/emccode/libstorage/api/types"
)

func GetVolumeID(libsClient types.Client, instanceId, service_id, plan_id string) (string, error) {
	return getScaleIOVolumeID(libsClient, instanceId)
}

func getScaleIOVolumeID(libsClient types.Client, instanceID string) (string, error) {
	volumeName, err := utils.CreateNameForScaleIO(instanceID)
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
