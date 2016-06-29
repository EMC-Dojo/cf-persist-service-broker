package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/EMC-Dojo/cf-persist-service-broker/model"
	"github.com/emccode/libstorage/api/types"
)

// ProjectDirectory : exports GOPATH for cf-
func ProjectDirectory() string {
	return filepath.Join(os.Getenv("GOPATH"), "src/github.com/EMC-Dojo/cf-persist-service-broker")
}

// CreatePlanIDString : Cloud Foundry, Open Source, The Way (To marshall Plan ID's). JUST LIKE MOM USED TO MAKE
func CreatePlanIDString(LibstorageService string) (string, error) {
	planID, err := json.Marshal(&model.PlanID{
		LibsHostName:    os.Getenv("LIBSTORAGE_URI"),
		LibsServiceName: LibstorageService,
	})
	if err != nil {
		return "", fmt.Errorf("Error creating PlanIDString : (%s)", err)
	}
	return string(planID[:len(planID)]), nil
}

// CreateVolumeRequest : generate a Libstorage volume request based on a provided name and storagepool as strings and size in GB
func CreateVolumeRequest(name, storagePool string, sizeInGB int64) (types.VolumeCreateRequest, error) {

	var az string
	var IOPS int64

	volumeName, err := CreateNameForVolume(name)

	volumeCreateRequest := types.VolumeCreateRequest{
		Name:             volumeName,
		AvailabilityZone: &az,       //AZ is not currently used with SIO / ISI
		IOPS:             &IOPS,     //IOPS is not currently used with SIO / ISI
		Size:             &sizeInGB, //Size is used, on SIO minimum size of 8GB, ISI does not use size
		Type:             &storagePool,
	}

	return volumeCreateRequest, err
}

// CreateNameForVolume : Follow constraints on Volume name based on old lame ScaleIO volume name problems #nochill
func CreateNameForVolume(instanceID string) (string, error) {
	if len(instanceID) < 32 {
		return instanceID, nil
	}

	cookedString := strings.Replace(instanceID, "-", "", -1)
	cookedString = cookedString[:len(cookedString)-1]

	if len(cookedString) > 31 {
		return "", fmt.Errorf("Volume name cannot exceed 32 characters when all hyphens are removed.")
	}

	return cookedString, nil
}
