package libstoragewrapper

import (
	"github.com/emccode/libstorage/api/context"
	"github.com/emccode/libstorage/api/types"
)

// CreateVolume : Creates Volume with Libstorage API
func CreateVolume(c types.APIClient, serviceName string, volumeCreateRequest types.VolumeCreateRequest) (*types.Volume, error) {
	return c.VolumeCreate(context.Background(), serviceName, &volumeCreateRequest)
}
