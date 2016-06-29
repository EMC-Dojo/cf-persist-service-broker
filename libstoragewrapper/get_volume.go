package libstoragewrapper

import (
	"fmt"

	"github.com/emccode/libstorage/api/context"

	"github.com/EMC-Dojo/cf-persist-service-broker/utils"
	"github.com/emccode/libstorage/api/types"
)

// GetVolumeAttachments : DOESN'T ACTUALLY RETURN ATTACHMENTS, NEEDS LIBSTORAGE PR
func GetVolumeAttachments(ctx types.Context, libsClient types.APIClient, serviceName, volumeID string) ([]*types.VolumeAttachment, error) {
	volume, err := libsClient.VolumeInspect(ctx, serviceName, volumeID, false)
	if err != nil {
		return nil, fmt.Errorf("error unable to get volume with ID %s. %s", volumeID, err)
	}

	return volume.Attachments, nil
}

// GetVolumeID : get the ID for a volume registered to serviceName in libstorage with name as instanceID
func GetVolumeID(libsClient types.APIClient, serviceName, instanceID string) (string, error) {
	parsedVolumeName, err := utils.CreateNameForVolume(instanceID)
	if err != nil {
		return "", fmt.Errorf("error parsing instanceID into volume name for InstanceID %s : (%s)", instanceID, err)
	}

	volumeID, err := getVolumeIDByName(libsClient, serviceName, parsedVolumeName)
	if err != nil {
		return "", fmt.Errorf("error getting volume id from volume name %s", err)
	}

	return volumeID, nil
}

func getVolumeIDByName(libsClient types.APIClient, serviceName, volumeName string) (string, error) {
	ctx := context.Background()
	volumes, err := libsClient.VolumesByService(ctx, serviceName, false)
	if err != nil {
		return "", fmt.Errorf("error getting volumes %s", err)
	}
	for _, v := range volumes {
		if v.Name == volumeName {
			return v.ID, nil
		}
	}

	return "", fmt.Errorf("could not find volume id from volume name: %s", volumeName)
}

// GetVolumeByID : query a service serviceName and fetch a volume by volumeID
func GetVolumeByID(libsClient types.APIClient, serviceName, volumeID string) (*types.Volume, error) {
	volume, err := libsClient.VolumeInspect(context.Background(), serviceName, volumeID, false)
	if err != nil {
		return nil, fmt.Errorf("error getting volume with ID %s : (%s)", volumeID, err)
	}
	return volume, nil
}
